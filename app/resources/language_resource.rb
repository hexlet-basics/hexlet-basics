class LanguageResource
  include Alba::Resource
  include Typelizer::DSL

  # root_key :user

  attributes :id, :slug, :members_count

  typelize :string, nullable: true
  attribute :duration do |language|
    language.duration
  end

  typelize :string, nullable: true
  attribute :name do |language|
    language.name
  end

  typelize :string, nullable: true
  attribute :cover do |language|
    "#{language.slug}.png"
  end
end
