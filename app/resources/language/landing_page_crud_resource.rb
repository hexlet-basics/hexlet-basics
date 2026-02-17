class Language::LandingPageCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Language::LandingPage

    typelize model: :string
    typelize relations: "Record<string, string>"
    typelize outcomes_image_thumb_url: [ :string, nullable: true ]
    typelize state_events: "Record<string, unknown>[]"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
    attribute(:outcomes_image_thumb_url) do
      urls = Rails.application.routes.url_helpers
      it.outcomes_image.attached? ? urls.rails_representation_url(it.outcomes_image.variant(:thumb)) : nil
    end
    attribute(:state_events) { it.class.enum_as_hashes(:states) }
  end

  typelize_from Language::LandingPage

  has_many :qna_items, resource: Language::LandingPageQnaItemCrudResource, key: "qna_items_attributes"
  has_one :language, resource: LanguageCrudResource
  has_one :landing_page_to_redirect, resource: Language::LandingPageForListsResource

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

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
