class AssistantChannel < ApplicationCable::Channel
  def subscribed
    lesson_member = Language::Lesson::Member.find(params[:id])
    stream_for lesson_member
    # stream_from "some_channel"
  end

  def unsubscribed
    # Any cleanup needed when channel is unsubscribed
  end
end
