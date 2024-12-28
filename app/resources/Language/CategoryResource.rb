class Language::CategoryResource
  include Alba::Resource
  include Typelizer::DSL

  # root_key :user

  attributes :id, :slug

  attribute :name do |category|
    category.name_ru
  end
end
