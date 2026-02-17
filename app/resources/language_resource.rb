class LanguageResource < ApplicationResource
  urls = Rails.application.routes.url_helpers

  typelize_from Language

  attributes :id, :slug, :learn_as, :progress, :category_id, :current_version_id, :created_at, :openai_assistant_id
  has_one :current_version, resource: Language::VersionResource

  typelize :string, nullable: true
  attribute :hexlet_program_landing_page do
    it.hexlet_program_landing_page.presence
  end

  typelize :string, nullable: true
  attribute :repository_url do
    it.repository_url
  end

  # typelize :string, nullable: true
  # attribute :main_landing_slug do |obj|
  #   obj.landing_pages.find_by(main: true)&.slug
  # end

  typelize :string
  attribute :cover_list_variant do
    urls.rails_representation_url(it.cover.variant(:list)) if it.cover.attached?
  end

  typelize :string
  attribute :cover_thumb_variant do
    urls.rails_representation_url(it.cover.variant(:thumb)) if it.cover.attached?
  end

  # TODO: хаха, вот я злодей. Реализовать нормальный подсчет рейтинга
  typelize :number
  attribute :rating_count do
    89
  end

  typelize :number
  attribute :members_count do
    it.members_count
  end

  typelize :number
  attribute :rating_value do
    it.id.even? ? 4 : 5
  end
end
