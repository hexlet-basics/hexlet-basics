class OriginalLanguageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language

  # root_key :user

  attributes :id, :name, :slug, :learn_as, :progress, :category_id

  typelize learn_as: [ enum: [ "first_language", "second_language" ] ]
  typelize progress: [ enum: [ "completed", "in_development", "draft" ] ]

  typelize :string, nullable: true
  attribute :repository_url do |obj|
    obj.repository_url
  end
end
