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
  attribute :cover_thumb_variant do |post|
    urls.rails_representation_url(post.cover.variant(:thumb)) if post.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_list_variant do |post|
    urls.rails_representation_url(post.cover.variant(:list)) if post.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_main_variant do |post|
    urls.rails_representation_url(post.cover.variant(:main)) if post.cover.attached?
  end

  typelize :number
  attribute :reading_time do |post|
    (post.body.split.size / 260).ceil
  end

  typelize :string
  attribute :url do |post|
    urls.blog_post_url(post.slug, suffix: I18n.locale == "en" ? "" : I18n.locale)
  end

  typelize :number
  attribute :likes_count do |post|
    post.id.digits.last
  end
end
