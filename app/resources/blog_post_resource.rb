class BlogPostResource
  include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  attributes :id, :name, :slug, :description, :body, :created_at, :state, :locale

  typelize :string, nullable: true
  attribute :cover_thumb_variant do |post|
    rails_representation_url(post.cover.variant(:thumb)) if post.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_list_variant do |post|
    rails_representation_url(post.cover.variant(:list)) if post.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_main_variant do |post|
    rails_representation_url(post.cover.variant(:main)) if post.cover.attached?
  end
end
