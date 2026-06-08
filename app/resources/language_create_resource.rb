class LanguageCreateResource < ApplicationResource
  typelize_from Language

  attributes :progress,
    :learn_as,
    :cover,
    :id,
    slug: [ String, true ],
    openai_assistant_id: [ String, true ],
    hexlet_program_landing_page: [ String, true ]

  typelize cover: "File | null"
  typelize id: [ :number, nullable: true ]
  typelize repository_url: [ :string, nullable: true ]
  attribute(:repository_url) { it.repository_url }
end
