# Language::Module::Version::Info
class Language::ModuleResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Module::Version::Info

  attributes :locale, :description, :name

  typelize :id, nullable: true
  attribute :id do |info|
    info.version.module.id
  end
end
