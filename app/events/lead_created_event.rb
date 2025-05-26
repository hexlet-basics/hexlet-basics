# typed: strict

class LeadCreatedEvent < TypedEvent
  DataShape = T.type_alias {
    {
      user_id: Integer,
      user_name: String,
      ym_client_id: T.nilable(String),
      email: String,
      phone: T.nilable(String),
      telegram: T.nilable(String),
      whatsapp: T.nilable(String),
      survey_answers_data: T::Array[T::Hash[String, T.untyped]],
      courses_data: T::Array[T::Hash[String, T.untyped]]
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
