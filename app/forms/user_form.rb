# frozen_string_literal: true

require 'reform/form/validation/unique_validator'

class UserForm < Reform::Form
  property :email
  property :password

  validates :email, presence: true, unique: true, format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :password, presence: true, length: { minimum: 6 }

  def email=(value)
    super(value.downcase)
  end
end
