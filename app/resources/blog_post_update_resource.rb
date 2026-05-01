class BlogPostUpdateResource < ApplicationResource
  typelize_from BlogPost

  attributes :id, :name, :slug, :description, :body, :state, :cover

  typelize name: :string
  typelize slug: :string
  typelize description: :string
  typelize body: :string
  typelize cover: "File | null"

  typelize cover_thumb_url: [ :string, nullable: true ]
  attribute(:cover_thumb_url) do
    urls = Rails.application.routes.url_helpers
    it.cover.attached? ? urls.rails_representation_url(it.cover.variant(:thumb)) : nil
  end
end
