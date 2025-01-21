class UserProfileFormResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User::ProfileForm

  attributes :first_name, :last_name

  typelize '"user_profile_form"', nullable: false
  attribute :type do |user|
    "user_profile_form"
  end
end
