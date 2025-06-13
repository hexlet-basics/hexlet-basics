class BlogPostCrudResource < ApplicationResource
  urls = Rails.application.routes.url_helpers

  typelize_from BlogPost
  root_key :blog_post

  attributes :id, :name, :slug, :description, :body, :state

  typelize_meta(meta: "{ cover_thumb_url: string, states: Record<string, unknown>[] }")
  meta do
    data = {
      cover_thumb_variant_url: object.cover.attached? ? urls.rails_representation_url(object.cover.variant(:thumb)) : nil,
      states: object.class.enum_as_hashes(:states)
    }
    data
  end
end
