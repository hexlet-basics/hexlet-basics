class LanguageResource
  include Alba::Resource
  include Typelizer::DSL

  # root_key :user

  attributes :id, :slug

  typelize :string, nullable: true
  attribute :name do |language|
    language.name
  end
end
