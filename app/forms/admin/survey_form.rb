class Admin::SurveyForm < Survey
  include ActiveFormModel

  permit :question,
    :description,
    :slug,
    :state,
    :parent_survey_item_id,
    :run_always,
    :run_after_finishing_lessons_count,
    items_attributes: [ :id, :value, :state, :order, :slug, :_destroy ]

  accepts_nested_attributes_for :items # , allow_destroy: true
end
