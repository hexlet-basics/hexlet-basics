class AssistantChannel < ApplicationCable::Channel
  def subscribed
    ai_chat = AiChat.find(params[:id])
    stream_for ai_chat
  end

  def unsubscribed
    # Any cleanup needed when channel is unsubscribed
  end
end
