# frozen_string_literal: true

class User::SignUpForm < User
  include ActiveFormModel

  permit :email, :password

  validates :password, presence: true, length: { minimum: 6 }

  def email=(email)
    if email.present?
      write_attribute(:email, email.downcase)
    else
      super
    end
  end
end
