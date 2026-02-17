class BlogPostResource < ApplicationResource
  urls = Rails.application.routes.url_helpers

  has_one :creator, resource: UserResource

  attributes :id,
    :name,
    :slug,
    :description,
    :body,
    :state,
    :locale,
    :created_at,
    :related_language_items_count

  typelize :string, nullable: true
  attribute :cover_thumb_variant do
    urls.rails_representation_url(it.cover.variant(:thumb)) if it.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_list_variant do
    urls.rails_representation_url(it.cover.variant(:list)) if it.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_main_variant do
    urls.rails_representation_url(it.cover.variant(:main)) if it.cover.attached?
  end

  typelize :number
  attribute :reading_time do
    (it.body.split.size / 260).ceil
  end

  typelize :string
  attribute :url do
    urls.blog_post_url(it.slug, suffix: I18n.locale == "en" ? "" : I18n.locale)
  end

  typelize :number
  attribute :likes_count do
    it.likes.count
  end
end
