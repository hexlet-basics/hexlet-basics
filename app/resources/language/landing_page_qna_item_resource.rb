class  Language::LandingPageQnaItemResource < ApplicationResource
  typelize_from Language::LandingPage::QnaItem

  attributes :id,
    :question,
    :answer

  typelize question: :string
  typelize answer: :string
end
