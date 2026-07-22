# typed: strict

class LanguageCreateResource < ApplicationResource
  typelize_from Language

  attributes :progress,
    :learn_as,
    :id,
    slug: [ String, true ],
    hexlet_program_landing_page: [ String, true ]

  # cover seeds the file input, not the stored attachment
  typelize cover: "File | null"
  attribute(:cover) { nil }
  typelize id: [ :number, nullable: true ]
end
