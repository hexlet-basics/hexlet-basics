class ReviewCreateResource < ApplicationResource
  typelize_from Review

  attributes :state,
    :pinned,
    id: [ Integer, true ],
    body: [ String, true ],
    first_name: [ String, true ],
    last_name: [ String, true ],
    language_id: [ Integer, true ],
    user_id: [ Integer, true ]
end
