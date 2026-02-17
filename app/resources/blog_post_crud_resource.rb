class BlogPostCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from BlogPost

    typelize model: :string
    typelize relations: "Record<string, string>"
    typelize states: "Record<string, unknown>[]"
    typelize cover_thumb_url: [ :string, nullable: true ]

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
    attribute(:states) { it.class.enum_as_hashes(:states) }
    attribute(:cover_thumb_url) do
      urls = Rails.application.routes.url_helpers
      it.cover.attached? ? urls.rails_representation_url(it.cover.variant(:thumb)) : nil
    end
  end

  typelize_from BlogPost

  attributes :id, :name, :slug, :description, :body, :state, :cover

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
