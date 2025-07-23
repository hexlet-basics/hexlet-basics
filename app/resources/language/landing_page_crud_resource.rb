class Language::LandingPageCrudResource < ApplicationResource
  urls = Rails.application.routes.url_helpers

  typelize_from Language::LandingPage
  root_key :data

  has_many :qna_items, resource: Language::LandingPageQnaItemCrudResource, key: "qna_items_attributes"
  has_one :language, resource: LanguageCrudResource
  has_one :landing_page_to_redirect, resource: Language::LandingPageCrudResource

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

  typelize_meta meta: "{ modelName: string, outcomes_image_thumb_url: string, state_events: Record<string, unknown>[] }"
  meta do
    {
      outcomes_image_thumb_url: object.outcomes_image.attached? ?
      urls.rails_representation_url(object.outcomes_image.variant(:thumb)) : nil,
      state_events: object.class.enum_as_hashes(:states),
      modelName: object.class.superclass.to_s.underscore
    }
  end
end
