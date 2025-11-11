class PagyResource < ApplicationResource
  attributes :count, :page, :last, :next, :from, :previous
  typelize count: :number, page: :number, last: :number, previous: :number
  typelize next: [ :number, nullable: true ]

  typelize :number
  attribute :first do
    1
  end
end
