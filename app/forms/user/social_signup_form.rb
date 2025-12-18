# frozen_string_literal: true

class User::SocialSignupForm < User
  include ActiveFormModel

  permit :email, :uid
end
