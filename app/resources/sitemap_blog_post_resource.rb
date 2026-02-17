class SitemapBlogPostResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from BlogPost

  attributes :id

  typelize :string
  attribute :name do
    it.name
  end

  typelize :string
  attribute :slug do
    it.slug
  end

  typelize "Locale"
  attribute :locale do
    it.locale
  end
end
