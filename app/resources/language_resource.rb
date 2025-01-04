class LanguageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Version::Info

  # root_key :user

  attributes :id, :description, :locale

  typelize :string, nullable: true
  attribute :duration do |info|
    info.language.duration
  end

  typelize :string, nullable: true
  attribute :slug do |info|
    info.language.slug
  end

  typelize :string, nullable: true
  attribute :name do |info|
    info.language.name
  end

  typelize :string, nullable: true
  attribute :cover do |info|
    "#{info.language.slug}.png"
  end

  # typelize :string, nullable: true
  # attribute :description do |language|
  #   language_version_info = language.current_version.infos.find_by!(locale: I18n.locale)
  #   language_version_info.description
  # end
end
