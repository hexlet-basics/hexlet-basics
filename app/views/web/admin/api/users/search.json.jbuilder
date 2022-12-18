# frozen_string_literal: true

json.array! @users do |user|
  json.id user.id
  json.text user.to_s
end
