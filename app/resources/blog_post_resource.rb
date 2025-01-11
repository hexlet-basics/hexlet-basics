class BlogPostResource
  include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  attributes :id, :name, :slug, :description, :body, :created_at, :state

  typelize :string, nullable: true
  attribute :cover_thumb_variant do |post|
    rails_representation_url(post.cover.variant(:thumb), only_path: true) if post.cover.attached?
  end

  typelize :string, nullable: true
  attribute :cover_list_variant do |post|
    rails_representation_url(post.cover.variant(:list), only_path: true) if post.cover.attached?
  end
end
