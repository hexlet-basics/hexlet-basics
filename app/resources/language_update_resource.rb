# typed: true

class LanguageUpdateResource < ApplicationResource
  typelize_from Language

  attributes :id, :progress, :learn_as, :slug, :openai_assistant_id, :hexlet_program_landing_page, :cover

  typelize slug: :string
  typelize cover: "File | null"

  typelize cover_thumb_url: [ :string, nullable: true ]
  attribute(:cover_thumb_url) do
    urls = Rails.application.routes.url_helpers
    it.cover.attached? ? urls.rails_representation_url(T.unsafe(it.cover).variant(:thumb)) : nil
  end

  typelize repository_url: [ :string, nullable: true ]
  attribute(:repository_url) { it.repository_url }
end
