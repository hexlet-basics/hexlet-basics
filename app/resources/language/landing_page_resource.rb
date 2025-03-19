class Language::LandingPageResource
  urls = Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::LandingPage

  # has_one :language, resource: LanguageResource

  attributes :id,
    :language_id,
    :slug,
    :main,
    :listed,
    :state,
    :order,
    :meta_title,
    :meta_description,
    :header,
    :description,
    :used_in_header,
    :used_in_description,
    :outcomes_header,
    :outcomes_description

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

  typelize :string, nullable: true
  attribute :outcomes_image do |lp|
    urls.rails_representation_url(lp.outcomes_image.variant(:main)) if lp.outcomes_image.attached?
  end

  # typelize :string
  # attribute :cover do |lp|
  #   "#{lp.language&.slug}.png"
  # end
end
