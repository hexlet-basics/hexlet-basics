# typed: strict

class SolutionCheckedEvent < TypedEvent
  DataShape = T.type_alias {
    {
      lesson_slug: String,
      course_slug: String,
      locale: Symbol,
      passed: T::Boolean
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
