# frozen_string_literal: true

class User::SocialSignupForm < User
  include ActiveFormModel

  permit :email, :uid

  def email=(email)
    if email.present?
      write_attribute(:email, email.downcase.strip)
    else
      super
    end
  end
end
