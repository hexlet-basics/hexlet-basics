# typed: strict

class Admin::LanguageLandingPageForm < Language::LandingPage
  include ActiveFormModel

  permit :meta_title,
    :meta_description,
    :listed,
    :slug,
    :footer,
    :footer_name,
    :order,
    :main,
    :state,
    :outcomes_image,
    :header,
    :name,
    :description,
    :language_id,
    :language_category_id,
    :landing_page_to_redirect_id,
    :used_in_header,
    :used_in_description,
    :outcomes_header,
    :outcomes_description

  skip_if_empty :outcomes_image
end
