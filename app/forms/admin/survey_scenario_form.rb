class Admin::SurveyScenarioForm < Survey::Scenario
  include ActiveFormModel

  permit :name,
    :state,
    :survey_item_id,
    # :event_name,
    # :event_treshold_count,
    items_attributes: [ :id, :survey_id, :order, :_destroy ],
    triggers_attributes: [ :id, :event_name, :event_threshold_count, :_destroy ]

  accepts_nested_attributes_for :items, allow_destroy: true
  accepts_nested_attributes_for :triggers, allow_destroy: true
end
