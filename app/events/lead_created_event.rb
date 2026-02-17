# typed: strict

class LeadCreatedEvent < ApplicationEvent
  NAME = "lead_created"

  DataShape = T.type_alias {
    {
      lead_id: Integer,
      user_id: Integer,
      user_name: String,
      first_name: T.nilable(String),
      last_name: T.nilable(String),
      ym_client_id: T.nilable(String),
      utm_source: T.nilable(String),
      utm_medium: T.nilable(String),
      utm_campaign: T.nilable(String),
      utm_term: T.nilable(String),
      utm_content: T.nilable(String),
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

  sig { returns(Integer) }
  def lead_id
    data.fetch(:lead_id)
  end
end
