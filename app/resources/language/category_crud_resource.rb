class Language::CategoryCrudResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Category
  root_key :language_category

  attributes :slug, :name, :header, :description

  meta do
    {}
  end
end
