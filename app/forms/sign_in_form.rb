# frozen_string_literal: true

class SignInForm < ApplicationForm
  property :email
  property :password

  validates :email, presence: true
  validates :password, presence: true
  validate :user_exists, :user_can_sign_in

  def user_can_sign_in
    if password.present? && !user&.valid_password?(password)
      errors.add(:password, :cannot_sign_in)
    end
  end

  def user_exists
    if email.present? && !user
      errors.add(:email, :user_does_not_exist)
    end
  end

  def user
    @user ||= User.find_by(email: email)
  end

  def email=(value)
    super(value.downcase)
  end
end
