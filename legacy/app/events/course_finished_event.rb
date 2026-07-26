# typed: strict

class CourseFinishedEvent < ApplicationEvent
  NAME = "course_finished"

  DataShape = T.type_alias {
    {
      occurrence_count: Integer,
      slug: String,
      locale: String
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
