class  Language::LandingPageQnaItemResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id,
    :question,
    :answer
end
