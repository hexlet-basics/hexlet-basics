# Language::Module::Version::Info
class Language::ModuleResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Module::Version::Info

  attributes :id, :locale, :description
 
  # typelize :string, nullable: true
  # attribute :description do |info|
  #   info.description
  # end
end
