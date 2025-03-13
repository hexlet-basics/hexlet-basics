class Language::LandingPageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::LandingPage

  # has_one :language, resource: LanguageResource

  attributes :id,
    :slug,
    :state,
    :order,
    :header,
    :meta_title,
    :meta_description,
    :description,
    :language_id

  typelize slug: :string
  typelize header: :string
  typelize meta_title: :string
  typelize meta_description: :string
  typelize description: :string

  typelize :number
  attribute :duration do |lp|
    I18n.t("common.hours", count: lp.language&.duration || 0)
  end

  typelize :number
  attribute :members_count do |lp|
    lp.language&.members_count || 0
  end

  # typelize :string
  # attribute :cover do |lp|
  #   "#{lp.language&.slug}.png"
  # end
end
