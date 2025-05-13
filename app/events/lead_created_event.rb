# typed: strict

class LeadCreatedEvent < TypedEvent
  DataShape = T.type_alias {
    {
      user_id: Integer,
      phone: String,
      telegram: String,
      whatsapp: String,
      survey_answers_data: T::Array[T::Hash[String, T.untyped]],
      courses_data: T::Array[T::Hash[String, T.untyped]]
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
