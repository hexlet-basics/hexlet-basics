# typed: strict

class LessonFinishedEvent < ApplicationEvent
  NAME = "lesson_finished"

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
