class Language::CategoryQnaItemResource < ApplicationResource
  typelize_from Language::Category::QnaItem

  attributes :id,
    :question,
    :answer
end
