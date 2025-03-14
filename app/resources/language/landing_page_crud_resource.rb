class Language::LandingPageCrudResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::LandingPage
  root_key :language_landing_page

  attributes :id,
    :slug,
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

  typelize_meta meta: "{ state_events: Array<[string, string]>}"
  meta do
    { state_events: object.aasm.events_for_select }
  end

  # typelize :state_events
end
