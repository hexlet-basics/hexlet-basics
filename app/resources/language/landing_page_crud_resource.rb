class Language::LandingPageCrudResource
  include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::LandingPage
  root_key :language_landing_page

  has_many :qna_items, resource: Language::LandingPageQnaItemCrudResource

  attributes :id,
    :slug,
    :footer,
    :footer_name,
    :state,
    :listed,
    :main,
    :order,
    :meta_title,
    :meta_description,
    :language_id,
    :language_category_id,
    :header,
    :description,
    :used_in_header,
    :used_in_description,
    :outcomes_header,
    :outcomes_description

  typelize_meta meta: "{ outcomes_image_thumb_url: string, state_events: Array<[string, string]>}"
  meta do
    {
      outcomes_image_thumb_url: object.outcomes_image.attached? ?
        rails_representation_url(object.outcomes_image.variant(:thumb)) : nil,
      state_events: object.aasm.events_for_select
    }
  end

  # typelize :state_events
end
