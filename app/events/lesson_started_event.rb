# typed: strict

class LessonStartedEvent < ApplicationEvent
  NAME = "lesson_started"

  DataShape = T.type_alias {
    {
      occurrence_count: Integer,
      lesson_slug: String,
      course_slug: String,
      locale: String
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
