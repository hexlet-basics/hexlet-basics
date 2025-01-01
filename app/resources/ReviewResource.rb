class ReviewResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id, :first_name, :last_name, :body
end
