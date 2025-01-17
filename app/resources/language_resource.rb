class LanguageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Version::Info

  # root_key :user

  attributes :id, :description, :locale, :created_at, :state, :order

  typelize :string, nullable: true
  attribute :duration do |info|
    I18n.t("common.hours", count: info.language.duration)
  end

  typelize :number, nullable: true
  attribute :id do |info|
    info.language_id
  end

  typelize :number, nullable: true
  attribute :members_count do |info|
    info.language.members_count
  end

  typelize :string, nullable: true
  attribute :state do |info|
    info.language.state
  end

  typelize :number, nullable: true
  attribute :order do |info|
    info.language.order
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
end
