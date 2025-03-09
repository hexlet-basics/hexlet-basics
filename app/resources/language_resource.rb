class LanguageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language

  attributes :id, :slug, :learn_as, :progress, :category_id, :current_version_id
  has_one :current_version, resource: Language::VersionResource

  typelize learn_as: [ enum: [ "first_language", "second_language" ] ]
  typelize progress: [ enum: [ "completed", "in_development", "draft" ] ]

  typelize :string, nullable: true
  attribute :repository_url do |obj|
    obj.repository_url
  end
end
