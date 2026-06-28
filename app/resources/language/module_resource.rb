# typed: strict

# Language::Module::Version::Info
class Language::ModuleResource < ApplicationResource
  typelize_from Language::Module::Version::Info

  attributes :locale, :description, :name

  typelize :number
  attribute :id do
    T.must(T.must(it.version).module).id
  end
end
