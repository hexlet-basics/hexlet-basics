# typed: strict

class BookRequestedEvent < ApplicationEvent
  NAME = "book_requested"

  DataShape = T.type_alias {
    {
      locale: String
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
