# Language::Module::Version
class Language::ModuleResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id
  #
  # typelize :string, nullable: true
  # attribute :name do |category|
  #   category.name_ru
  # end
end
