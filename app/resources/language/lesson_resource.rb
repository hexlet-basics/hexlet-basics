class Language::LessonResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id, :slug, :module_id

  # typelize :number, nullable: true
  # attribute :module_id do |category|
  #   category.name_ru
  # end
end
