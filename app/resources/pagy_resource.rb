class PagyResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :count, :page, :last, :prev, :next, :from
  typelize count: :number, page: :number, last: :number
end
