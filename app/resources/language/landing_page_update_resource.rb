class Language::LandingPageUpdateResource < ApplicationResource
  typelize_from Language::LandingPage

  attributes :id,
    :slug,
    :footer,
    :footer_name,
    :state,
    :listed,
    :main,
    :name,
    :order,
    :meta_title,
    :meta_description,
    :language_id,
    :landing_page_to_redirect_id,
    :header,
    :description,
    :used_in_header,
    :used_in_description,
    :outcomes_header,
    :outcomes_image,
    :outcomes_description

  typelize slug: :string
  typelize name: :string
  typelize meta_title: :string
  typelize header: :string
  typelize outcomes_image: "File | null"

  typelize outcomes_image_thumb_url: [ :string, nullable: true ]
  attribute(:outcomes_image_thumb_url) do
    urls = Rails.application.routes.url_helpers
    it.outcomes_image.attached? ? urls.rails_representation_url(it.outcomes_image.variant(:thumb)) : nil
  end
end
