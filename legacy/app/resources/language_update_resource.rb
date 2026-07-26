# typed: strict

class LanguageUpdateResource < ApplicationResource
  typelize_from Language

  attributes :id, :progress, :learn_as, :slug, :hexlet_program_landing_page

  typelize slug: :string

  # cover seeds the file input, not the stored attachment; the current one is exposed via cover_thumb_url
  typelize cover: "File | null"
  attribute(:cover) { nil }

  typelize cover_thumb_url: [ :string, nullable: true ]
  attribute(:cover_thumb_url) do
    urls = Rails.application.routes.url_helpers
    it.cover.attached? ? urls.rails_representation_url(T.unsafe(it.cover).variant(:thumb)) : nil
  end

  typelize repository_url: [ :string, nullable: true ]
  attribute(:repository_url) { it.repository_url }
end
