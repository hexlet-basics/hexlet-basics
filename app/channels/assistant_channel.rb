# typed: strict

class AssistantChannel < ApplicationCable::Channel
  extend T::Sig

  sig { void }
  def subscribed
    ai_chat = AiChat.find(params[:id])
    stream_for ai_chat
  end

  sig { void }
  def unsubscribed
    # Any cleanup needed when channel is unsubscribed
  end
end
