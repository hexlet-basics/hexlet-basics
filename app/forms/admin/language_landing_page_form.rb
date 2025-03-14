class Admin::LanguageLandingPageForm < Language::LandingPage
  include ActiveFormModel

  permit :meta_title,
    :meta_description,
    :listed,
    :slug,
    :order,
    :main,
    :state_event,
    :outcomes_image,
    :header,
    :description,
    :language_id,
    :language_category_id,
    :used_in_header,
    :used_in_description,
    :outcomes_header,
    :outcomes_description,
    qna_items_attributes: [ :id, :question, :answer, :_destroy ]

  accepts_nested_attributes_for :qna_items, allow_destroy: true
end
