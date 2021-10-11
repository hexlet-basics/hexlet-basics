# frozen_string_literal: true

class User::SocialSignupForm < User
  include ActiveFormModel

  fields :email, :uid

  def email=(email)
    if email.present?
      write_attribute(:email, email.downcase)
    else
      super
    end
  end
end
