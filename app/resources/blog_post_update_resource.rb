# typed: strict

class BlogPostUpdateResource < ApplicationResource
  typelize_from BlogPost

  attributes :id, :name, :slug, :description, :state, :cover

  typelize name: :string
  typelize slug: :string
  typelize description: :string
  typelize cover: "File | null"

  typelize :string
  attribute :rich_body do
    next "" unless it.rich_body.body.present?

    BlogPostRichTextContent.to_editor_html(it.rich_body.body.to_html)
  end

  typelize cover_thumb_url: [ :string, nullable: true ]
  attribute(:cover_thumb_url) do
    urls = Rails.application.routes.url_helpers
    it.cover.attached? ? urls.rails_representation_url(T.unsafe(it.cover).variant(:thumb)) : nil
  end
end
