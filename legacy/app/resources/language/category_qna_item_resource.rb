# typed: strict

class Language::CategoryQnaItemResource < ApplicationResource
  typelize_from Language::Category::QnaItem

  attributes :id,
    :question,
    :answer

  typelize question: :string
  typelize answer: :string

  typelize :boolean
  attribute :_destroy do
    false
  end
end
