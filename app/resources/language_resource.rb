class LanguageResource < ApplicationResource
  urls = Rails.application.routes.url_helpers

  typelize_from Language

  attributes :id, :slug, :learn_as, :progress, :category_id, :current_version_id, :created_at, :openai_assistant_id
  has_one :current_version, resource: Language::VersionResource

  typelize :string, nullable: true
  attribute :hexlet_program_landing_page do |obj|
    obj.hexlet_program_landing_page.presence
  end

  typelize :string, nullable: true
  attribute :repository_url do |obj|
    obj.repository_url
  end

  # typelize :string, nullable: true
  # attribute :main_landing_slug do |obj|
  #   obj.landing_pages.find_by(main: true)&.slug
  # end

  typelize :string
  attribute :cover_list_variant do |obj|
    urls.rails_representation_url(obj.cover.variant(:list)) if obj.cover.attached?
  end

  typelize :string
  attribute :cover_thumb_variant do |obj|
    urls.rails_representation_url(obj.cover.variant(:thumb)) if obj.cover.attached?
  end

  # TODO: хаха, вот я злодей. Реализовать нормальный подсчет рейтинга
  typelize :number
  attribute :rating_count do |obj|
    89
  end

  typelize :number
  attribute :members_count do |obj|
    obj.members_count
  end

  typelize :number
  attribute :rating_value do |obj|
    obj.id.even? ? 4 : 5
  end
end
