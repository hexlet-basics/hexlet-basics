# typed: strict

class UserSignedUpEvent < ApplicationEvent
  NAME = "user_signed_up"

  DataShape = T.type_alias do
    {
      user_id: Integer,
      email: String,
      first_name: T.nilable(String),
      last_name: T.nilable(String),
      locale: Symbol
    }
  end

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
