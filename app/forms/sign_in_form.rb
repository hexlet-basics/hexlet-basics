# frozen_string_literal: true

class SignInForm < ApplicationForm
  property :email
  property :password

  validates :email, presence: true
  validates :password, presence: true
  validate :user_can_sign_in

  def user_can_sign_in
    errors.add(:password, :cannot_sign_in) unless user&.authenticate(password)
  end

  def user
    @user ||= User.find_by(email: email)
  end

  def email=(value)
    super(value.downcase)
  end
end
