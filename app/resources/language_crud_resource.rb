class LanguageCrudResource
  include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language
  root_key :language

  attributes :id, :progress, :learn_as, :slug, :openai_assistant_id

  typelize_meta meta: "{ cover_signed_id: string, cover_thumb_url: string, repository_url: string, slug: string }"
  meta do
    {
      # cover_signed_id: object.cover.signed_id,
      # state_events: object.aasm.events_for_select,
      cover_thumb_variant_url: object.cover.attached? ?
        rails_representation_url(object.cover.variant(:thumb)) : nil,
      repository_url: object.repository_url,
      slug: object.slug
    }
  end
end
