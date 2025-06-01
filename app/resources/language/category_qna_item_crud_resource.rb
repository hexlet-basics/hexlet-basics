class Language::CategoryQnaItemCrudResource < ApplicationResource
  typelize_from Language::Category::QnaItem

  attributes :id,
    :question,
    :answer

  meta do
    {}
  end
end
