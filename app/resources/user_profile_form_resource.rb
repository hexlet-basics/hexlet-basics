class UserProfileFormResource < ApplicationResource
  typelize_from User::ProfileForm

  attributes :first_name, :last_name, :contact_method, :contact_value

  typelize '"user_profile_form"', nullable: false
  attribute :type do |user|
    "user_profile_form"
  end

  meta do
    {}
  end
end
