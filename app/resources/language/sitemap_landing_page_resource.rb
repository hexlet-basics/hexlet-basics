class Language::SitemapLandingPageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::LandingPage

  attributes :id, :language_id, :slug, :header

  typelize slug: :string
  typelize header: :string

  typelize "Locale"
  attribute :locale do |info|
    info.locale
  end
end
