class Language::SitemapLandingPageResource < ApplicationResource
  typelize_from Language::LandingPage

  attributes :id, :language_id, :slug, :header

  typelize slug: :string
  typelize header: :string

  typelize "Locale"
  attribute :locale do
    it.locale
  end
end
