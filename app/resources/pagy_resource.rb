class PagyResource < ApplicationResource
  attributes :count, :page, :last, :prev, :next, :from
  typelize count: :number, page: :number, last: :number
end
