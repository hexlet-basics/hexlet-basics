# frozen_string_literal: true

json.array! @users do |user|
  json.id user.id
  json.text "#{user.email} (#{user.first_name} #{user.last_name})"
end
