# typed: strict

class UserSignedInEvent < ApplicationEvent
  NAME = "user_signed_in"

  DataShape = T.type_alias do
    {
      user_id: Integer,
      occurrence_count: Integer,
      email: String,
      locale: Symbol
    }
  end

  sig do
    params(
      data: DataShape,
      event_id: T.nilable(String),
      metadata: ApplicationEvent::EventMetadata
    ).void
  end
  def initialize(data:, event_id: SecureRandom.uuid, metadata: {})
    super
  end
end
