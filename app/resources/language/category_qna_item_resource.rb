class Language::CategoryQnaItemResource < ApplicationResource
  typelize_from Language::Category::QnaItem

  attributes :id,
    :question,
    :answer

  typelize :boolean
  attribute :_destroy do
    false
  end
end
