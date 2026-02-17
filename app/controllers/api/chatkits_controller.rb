# frozen_string_literal: true

class Api::ChatkitsController < Api::ApplicationController
  include ActionController::Live

  STREAMING_TYPES = %w[
    threads.create
    threads.add_user_message
    threads.add_client_tool_output
    threads.retry_after_item
    threads.custom_action
  ].freeze

  before_action :require_api_auth!

  def create
    payload = parse_payload

    if STREAMING_TYPES.include?(payload["type"])
      process_streaming(payload)
    else
      render json: process_non_streaming(payload)
    end
  rescue ActiveRecord::RecordNotFound
    render json: { error: "not_found" }, status: :not_found
  rescue JSON::ParserError
    render json: { error: "invalid_json" }, status: :unprocessable_entity
  rescue StandardError => e
    render json: { error: e.message }, status: :unprocessable_entity
  end

  private

  def parse_payload
    return params.to_unsafe_h if params[:type]

    JSON.parse(request.raw_post)
  end

  def process_non_streaming(payload)
    lesson_member = find_lesson_member!(payload)
    thread_id = lesson_member.openai_thread_id

    case payload["type"]
    when "threads.list"
      data = []
      if thread_id.present?
        data << serialize_thread(lesson_member:, include_items: false)
      end
      { data:, has_more: false, after: nil }
    when "threads.get_by_id"
      requested_thread_id = payload.dig("params", "thread_id")
      raise ActiveRecord::RecordNotFound if requested_thread_id != thread_id

      serialize_thread(lesson_member:, include_items: true)
    when "items.list"
      requested_thread_id = payload.dig("params", "thread_id")
      raise ActiveRecord::RecordNotFound if requested_thread_id != thread_id

      serialize_items_page(
        lesson_member:,
        order: payload.dig("params", "order"),
        limit: payload.dig("params", "limit"),
        after: payload.dig("params", "after"),
      )
    when "threads.update"
      serialize_thread(lesson_member:, include_items: false)
    when "threads.delete"
      lesson_member.update!(openai_thread_id: nil)
      lesson_member.messages.delete_all
      {}
    when "items.feedback", "attachments.create", "attachments.delete", "input.transcribe"
      {}
    else
      raise "Unsupported ChatKit request type: #{payload["type"]}"
    end
  end

  def process_streaming(payload)
    response.headers["Content-Type"] = "text/event-stream"
    response.headers["Cache-Control"] = "no-cache"
    response.headers["X-Accel-Buffering"] = "no"
    response.headers["Last-Modified"] = Time.current.httpdate

    lesson_member = find_lesson_member!(payload)
    can_create_assistant_message = Language::Lesson::Member::MessagePolicy.new(
      current_user,
      Language::Lesson::Member::Message
    ).create?

    unless can_create_assistant_message
      write_sse(
        type: "notice",
        level: "warning",
        message: I18n.t(
          "languages.lessons.show.chat.disabled_html",
          scope: "web",
          default: "Assistant messages are temporarily unavailable."
        ),
      )
      return
    end

    write_sse(type: "stream_options", stream_options: { allow_cancel: true })

    case payload["type"]
    when "threads.create"
      ensure_openai_thread!(lesson_member)
      write_sse(
        type: "thread.created",
        thread: serialize_thread(lesson_member:, include_items: false),
      )
      handle_user_message_stream(lesson_member:, payload:)
    when "threads.add_user_message"
      ensure_openai_thread!(lesson_member)
      handle_user_message_stream(lesson_member:, payload:)
    else
      write_sse(
        type: "error",
        code: "custom",
        message: "Unsupported streaming request: #{payload["type"]}",
        allow_retry: false,
      )
    end
  rescue StandardError => e
    write_sse(
      type: "error",
      code: "custom",
      message: e.message,
      allow_retry: true,
    )
  ensure
    response.stream.close
  end

  def handle_user_message_stream(lesson_member:, payload:)
    input = payload.dig("params", "input") || {}
    metadata = payload["metadata"] || {}
    thread_id = lesson_member.openai_thread_id
    user_text = extract_user_text(input)

    user_message = create_db_message(
      lesson_member:,
      role: "user",
      body: user_text,
    )

    write_sse(type: "thread.item.done", item: serialize_message_item(user_message, thread_id:))

    assistant_item_id = "msg_#{SecureRandom.hex(8)}"
    assistant_content = { type: "output_text", text: "", annotations: [] }

    write_sse(
      type: "thread.item.added",
      item: {
        id: assistant_item_id,
        type: "assistant_message",
        thread_id:,
        created_at: Time.current.iso8601,
        content: []
      },
    )
    write_sse(
      type: "thread.item.updated",
      item_id: assistant_item_id,
      update: {
        type: "assistant_message.content_part.added",
        content_index: 0,
        content: assistant_content
      },
    )

    assistant_text = +""

    stream_assistant_response(
      lesson_member:,
      message: user_text,
      user_code: metadata["user_code"].to_s,
      output: metadata["output"].to_s,
    ) do |delta|
      assistant_text << delta
      write_sse(
        type: "thread.item.updated",
        item_id: assistant_item_id,
        update: {
          type: "assistant_message.content_part.text_delta",
          content_index: 0,
          delta:
        },
      )
    end

    final_content = {
      type: "output_text",
      text: assistant_text,
      annotations: []
    }

    write_sse(
      type: "thread.item.updated",
      item_id: assistant_item_id,
      update: {
        type: "assistant_message.content_part.done",
        content_index: 0,
        content: final_content
      },
    )

    assistant_message = create_db_message(
      lesson_member:,
      role: "assistant",
      body: assistant_text,
    )

    write_sse(
      type: "thread.item.done",
      item: serialize_message_item(assistant_message, thread_id:),
    )
  end

  def stream_assistant_response(lesson_member:, message:, user_code:, output:)
    lesson = lesson_member.lesson
    lesson_info = lesson.infos.find_by!(locale: I18n.locale)
    language = lesson.language

    unless language.openai_assistant_id
      raise "#{language} without openai_assistant_id"
    end

    openai_api = DepsLocator.current.openai_api

    openai_api.messages.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        role: "user",
        content: [
          {
            type: "text",
            text: "Пользовательский код (решение практики): #{user_code}"
          },
          {
            type: "text",
            text: "Результат запуска пользовательского кода (используя тесты задания): #{output}"
          },
          {
            type: "text",
            text: "Вопрос пользователя: #{message}"
          }
        ]
      }
    )

    instructions = "
      Ты помогаешь изучать #{language.slug} на основе загруженного курса в файлах.
      Этот тред посвящен уроку #{lesson_info.name}.
      Курс состоит из теории, практики и тестов, которые выполняются прямо в браузере.
      Ты не показываешь решение практики, пользователь должен решить практику самостоятельно.
      Направляй, давай объяснения, помогай разобраться, предлагай шаги для решения, выдвигай гипотезы.
      Отвечай на языке: #{I18n.t(I18n.locale, scope: 'common.languages')}. Отвечай коротко.

      Если во время обсуждения обнаружится ошибка в теории или задании (описании, решении, тестах),
      то порекомендуй написать об этом в сообществе https://t.me/HexletLearningBot. Там сидит команда проекта.
    "

    openai_api.runs.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        assistant_id: language.openai_assistant_id,
        instructions:,
        stream: proc do |chunk|
          next unless chunk["object"] == "thread.message.delta"

          content = chunk.dig("delta", "content") || []
          texts = content.map { it.dig("text", "value") }.compact
          next if texts.blank?

          texts.each { |delta| yield delta }
        end
      }
    )
  end

  def ensure_openai_thread!(lesson_member)
    return if lesson_member.openai_thread_id.present?

    thread = DepsLocator.current.openai_api.threads.create
    lesson_member.update!(openai_thread_id: thread["id"])
  end

  def create_db_message(lesson_member:, role:, body:)
    lesson = lesson_member.lesson

    lesson_member.messages.create!(
      body:,
      role:,
      language_lesson: lesson,
      language: lesson.language,
      user: lesson_member.user,
    )
  end

  def serialize_thread(lesson_member:, include_items:)
    {
      id: lesson_member.openai_thread_id,
      title: nil,
      created_at: lesson_member.created_at.iso8601,
      status: { type: "active" },
      metadata: {},
      items: include_items ? serialize_items_page(lesson_member:) : empty_page
    }
  end

  def serialize_items_page(lesson_member:, order: "desc", limit: 20, after: nil)
    relation = lesson_member.messages.order(created_at: normalize_order(order))
    relation = apply_after_cursor(relation:, after:, order:)
    limit_value = normalize_limit(limit)
    scoped_relation = relation.limit(limit_value)
    messages = scoped_relation.to_a
    has_more = relation.offset(limit_value).exists?
    last_id = messages.last&.id

    {
      data: messages.map { |message| serialize_message_item(message, thread_id: lesson_member.openai_thread_id) },
      has_more:,
      after: has_more && last_id ? "dbmsg_#{last_id}" : nil
    }
  end

  def serialize_message_item(message, thread_id:)
    if message.role_user?
      {
        id: "dbmsg_#{message.id}",
        type: "user_message",
        thread_id:,
        created_at: message.created_at.iso8601,
        content: [
          {
            type: "input_text",
            text: message.body.to_s
          }
        ],
        attachments: [],
        quoted_text: nil,
        inference_options: {}
      }
    else
      {
        id: "dbmsg_#{message.id}",
        type: "assistant_message",
        thread_id:,
        created_at: message.created_at.iso8601,
        content: [
          {
            type: "output_text",
            text: message.body.to_s,
            annotations: []
          }
        ]
      }
    end
  end

  def apply_after_cursor(relation:, after:, order:)
    return relation if after.blank?

    after_id = after.to_s.delete_prefix("dbmsg_").to_i
    return relation if after_id.zero?

    if normalize_order(order) == :desc
      relation.where("id < ?", after_id)
    else
      relation.where("id > ?", after_id)
    end
  end

  def normalize_order(order)
    order.to_s == "asc" ? :asc : :desc
  end

  def normalize_limit(limit)
    parsed = Integer(limit || 20)
    return 1 if parsed <= 0
    return 100 if parsed > 100

    parsed
  rescue ArgumentError, TypeError
    20
  end

  def empty_page
    { data: [], has_more: false, after: nil }
  end

  def extract_user_text(input)
    content = Array(input["content"])
    text = content.filter_map do |part|
      next unless part["type"] == "input_text"

      part["text"].to_s
    end.join("\n").strip

    text.presence || ""
  end

  def find_lesson_member!(payload)
    lesson_id = payload.dig("metadata", "lesson_id")
    raise ActiveRecord::RecordNotFound if lesson_id.blank?

    lesson = Language::Lesson.find(lesson_id)
    lesson.members.find_by!(user: current_user)
  end

  def write_sse(payload)
    response.stream.write("data: #{JSON.generate(payload)}\n\n")
  end
end
