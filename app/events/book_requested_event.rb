# typed: strict

class BookRequestedEvent < TypedEvent
  # DataShape = T.type_alias {
  #   {
  #     slug: String,
  #     locale: Symbol
  #   }
  # }
  #
  # sig { params(data: DataShape, kwargs: T.untyped).void }
  # def initialize(data:, **kwargs)
  #   super
  # end
end
