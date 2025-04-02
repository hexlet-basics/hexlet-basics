class Language::CategoryCrudResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Category
  root_key :language_category

  attributes :id, :slug, :name

  meta do
    {}
  end
end
