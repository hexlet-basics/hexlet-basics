class Language::CategoryResource < ApplicationResource
  typelize_from Language::Category
  # root_key :user

  attributes :id, :slug, :name, :created_at, :header, :description
  # typelize name: :string, nullable: true
  typelize name: :string
  # attribute :name do |category|
  #   category.name_ru
  # end

  typelize "Locale"
  attribute :locale do
    it.locale
  end
end
