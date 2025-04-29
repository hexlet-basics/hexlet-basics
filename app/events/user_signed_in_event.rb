# typed: strict

class UserSignedInEvent < TypedEvent
  DataShape = T.type_alias do
    {
      id: Integer,
      email: String
    }
  end

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
