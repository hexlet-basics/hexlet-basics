class OriginalLanguageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language

  # root_key :user

  attributes :id, :name, :slug, :learn_as, :state
end
