# typed: strict

class Language::LandingPageQnaItemCrudResource < ApplicationResource
  typelize_from Language::LandingPage::QnaItem
  # root_key :data

  attributes :id,
    :question,
    :answer

  typelize id: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do
    false
  end
end
