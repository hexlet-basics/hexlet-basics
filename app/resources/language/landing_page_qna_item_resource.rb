class  Language::LandingPageQnaItemResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::LandingPage::QnaItem

  attributes :id,
    :question,
    :answer

  typelize question: :string
  typelize answer: :string
end
