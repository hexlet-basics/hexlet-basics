# typed: strict

class LessonStartedEvent < TypedEvent
  DataShape = T.type_alias {
    {
      occurrence_count: Integer,
      lesson_slug: String,
      course_slug: String,
      locale: Symbol
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
