class Admin::LanguageLandingPageForm < Language::LandingPage
  include ActiveFormModel

  permit :meta_title,
    :meta_description,
    :header,
    :state_event,
    :description,
    :language_id,
    :language_category_id,
    :slug,
    :order
end
