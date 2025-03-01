class SitemapBlogPostResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from BlogPost

  attributes :id

  typelize :string
  attribute :name do |post|
    post.name
  end

  typelize :string
  attribute :slug do |post|
    post.slug
  end

  typelize "Locale"
  attribute :locale do |post|
    post.locale
  end
end
