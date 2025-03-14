class Admin::LanguageLandingPageForm < Language::LandingPage
  include ActiveFormModel

  permit :meta_title,
    :meta_description,
    :listed,
    :slug,
    :order,
    :main,
    :state_event,
    :header,
    :description,
    :language_id,
    :language_category_id,
    :used_in_header,
    :used_in_description,
    :outcomes_header,
    :outcomes_description
end
