class Language::CategoryResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Category
  # root_key :user

  attributes :id, :slug, :name, :created_at
  # typelize name: :string, nullable: true
  typelize name: :string
  # attribute :name do |category|
  #   category.name_ru
  # end

  typelize "Locale"
  attribute :locale do |category|
    category.locale
  end
end
