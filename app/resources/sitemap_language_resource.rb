class SitemapLanguageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Version::Info

  typelize "Locale"
  attribute :locale do |info|
    info.locale
  end

  typelize :number
  attribute :id do |info|
    info.language_id
  end

  typelize :string
  attribute :slug do |info|
    info.language.slug
  end

  typelize :string
  attribute :name do |info|
    info.header
  end
end
