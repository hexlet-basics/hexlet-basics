# typed: true

class ReviewUpdateResource < ApplicationResource
  typelize_from Review

  attributes :id,
    :state,
    :language_id,
    :user_id,
    :pinned,
    :last_name,
    body: String,
    first_name: String
end
