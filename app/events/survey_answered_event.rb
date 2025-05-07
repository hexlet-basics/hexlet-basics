# typed: strict

class SurveyAnsweredEvent < TypedEvent
  DataShape = T.type_alias {
    {
      survey_answer_id: Integer,
      survey_scenario_member_id: Integer,
      next_survey_id: T.nilable(Integer)
    }
  }

  sig { params(data: DataShape, kwargs: T.untyped).void }
  def initialize(data:, **kwargs)
    super
  end
end
