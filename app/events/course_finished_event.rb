# typed: strict

class CourseFinishedEvent < TypedEvent
  DataShape = T.type_alias {
    {
      occurrence_count: Integer,
      slug: String,
      locale: Symbol
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
