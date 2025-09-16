class PagyResource < ApplicationResource
  attributes :count, :page, :last, :next, :from
  typelize count: :number, page: :number, last: :number
  typelize next: [ :number, nullable: true ]

  typelize :number, nullable: true
  attribute :previous do |pagy|
    pagy.prev
  end

  typelize :number
  attribute :first do
    1
  end
end
