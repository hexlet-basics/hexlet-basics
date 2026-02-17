# typed: strict

class SurveyStartedEvent < ApplicationEvent
  NAME = "survey_started"

  DataShape = T.type_alias {
    {
      user_id: Integer,
      email: String,
      occurrence_count: Integer,
      survey_id: Integer,
      survey_scenario_member_id: Integer,
      locale: Symbol
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
