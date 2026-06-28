# typed: true

class ReviewCreateResource < ApplicationResource
  typelize_from Review

  attributes :state,
    :pinned,
    :id,
    :language_id,
    :user_id,
    body: [ String, true ],
    first_name: [ String, true ],
    last_name: [ String, true ]

  typelize id: [ :number, nullable: true ]
  typelize language_id: [ :number, nullable: true ]
  typelize user_id: [ :number, nullable: true ]
end
