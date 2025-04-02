class Language::LandingPageQnaItemCrudResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id,
    :question,
    :answer

  meta do
    {}
  end
end
