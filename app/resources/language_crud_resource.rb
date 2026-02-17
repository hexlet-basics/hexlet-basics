class LanguageCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Language

    typelize model: :string
    typelize relations: "Record<string, string>"
    typelize cover_thumb_url: [ :string, nullable: true ]
    typelize repository_url: [ :string, nullable: true ]
    typelize slug: [ :string, nullable: true ]

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
    attribute(:cover_thumb_url) do
      urls = Rails.application.routes.url_helpers
      it.cover.attached? ? urls.rails_representation_url(it.cover.variant(:thumb)) : nil
    end
    attribute(:repository_url) { it.repository_url }
    attribute(:slug) { it.slug }
  end

  typelize_from Language

  attributes :id, :progress, :learn_as, :slug, :openai_assistant_id, :hexlet_program_landing_page, :cover

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
