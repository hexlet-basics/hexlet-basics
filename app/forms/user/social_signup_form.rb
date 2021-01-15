# frozen_string_literal: true

class User::SocialSignupForm < User
  include ActiveFormModel

  permit :email, :uid

  def email=(value)
    if email.present?
      write_attribute(:email, email.downcase)
    else
      super
    end
  end
end
