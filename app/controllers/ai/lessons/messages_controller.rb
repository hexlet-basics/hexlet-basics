# typed: strict

class Ai::Lessons::MessagesController < Ai::ApplicationController
  sig { void }
  def index
    lesson = Language::Lesson.find(params.expect(:lesson_id))
    lesson_member = lesson.members.find_by!(user: current_user)
    ai_chat = lesson_member.ai_chat

    messages = ai_chat ? AiMessageResource.new(ai_chat.ai_messages.where(role: %w[user assistant]).order(:id)) : []

    render json: messages
  end

  sig { void }
  def create
    unless AiMessagePolicy.new(current_user, AiMessage).create?
      head :too_many_requests
      return
    end

    lesson = Language::Lesson.find(params.expect(:lesson_id))
    lesson_member = lesson.members.find_by!(user: current_user)
    ai_chat = AiChat.find_or_create_by!(user: lesson_member.user, language_lesson_member: lesson_member)

    Assistants::RunJob.perform_later(
      ai_chat_id: ai_chat.id,
      message: params[:message],
      output: params[:output],
      user_code: params[:user_code],
      locale: I18n.locale.to_s
    )

    head :created
  end
end
