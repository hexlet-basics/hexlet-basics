class Ai::Lessons::MessagesController < Ai::ApplicationController
  # include ActionController::Live
  # include Import["openai_api"]

  def index
    lesson = Language::Lesson.find(params[:lesson_id])
    # language = lesson_info.language
    lesson_member = lesson.members.find_by!(user: current_user)

    messages = []
    if lesson_member.openai_thread_id?
      openai_api = OpenAI::Client.new do |f|
        if configus.hexlet_proxy.url.present?
          f.proxy = { uri: configus.hexlet_proxy.url }
        end
      end
      result = openai_api.messages.list(thread_id: lesson_member.openai_thread_id)
      messages = result["data"].map do
        {
          id: it["id"],
          role: it["role"],
          created_at: it["created_at"],
          content: it["content"].map { it["text"]["value"] }.join
        }
      end.reverse
    end

    render json: messages
  end

  def create
    lesson = Language::Lesson.find(params[:lesson_id])
    lesson_member = lesson.members.find_by!(user: current_user)

    m = lesson_member.messages.build body: params[:message]
    m.role = "user"
    m.language_lesson = lesson
    m.language = lesson.language
    m.save!

    Assistants::RunJob.perform_later(
      lesson_member_id: lesson_member.id,
      message: params[:message],
      output: params[:output],
      user_code: params[:user_code]
    )

    head :created
  end
end
