class Language::LandingPageForListsResource < ApplicationResource
  # include Rails.application.routes.url_helpers

  typelize_from Language::LandingPage

  has_one :language # , resource: LanguageResource
  attributes :id, :slug, :header, :language_id, :locale, :name

  typelize slug: :string
  typelize header: :string
  typelize language_id: :string

  typelize :number
  attribute :duration do |lp|
    I18n.t("common.hours", count: lp.language&.duration || 0)
  end

  typelize :number
  attribute :members_count do |lp|
    lp.language&.members_count || 0
  end
end
