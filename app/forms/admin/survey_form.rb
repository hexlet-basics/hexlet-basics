class Admin::SurveyForm < Survey
  include ActiveFormModel

  permit :question,
    :description,
    :slug,
    :state,
    items_attributes: [ :id, :value, :state, :order, :_destroy ]

  accepts_nested_attributes_for :items, allow_destroy: true
end
