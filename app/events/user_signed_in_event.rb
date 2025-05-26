# typed: strict

class UserSignedInEvent < TypedEvent
  DataShape = T.type_alias do
    {
      id: Integer,
      email: String,
      first_name: T.nilable(String),
      last_name: T.nilable(String)
    }
  end

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
