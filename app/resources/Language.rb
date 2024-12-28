class LanguageResource
  include Alba::Resource
  include Typelizer::DSL

  # root_key :user

  attributes :id, :slug
end

