# frozen_string_literal: true

class SignInForm
  include ActiveFormModel::Virtual

  fields :email, :password

  validates :email, presence: true
  validates :password, presence: true
  validate :user_exists, :user_can_sign_in

  def user_can_sign_in
    errors.add(:password, :cannot_sign_in) if password.present? && !user&.valid_password?(password)
  end

  def user_exists
    errors.add(:email, :user_does_not_exist_html) if email.present? && !user
  end

  def user
    @user ||= User.active.find_by(email: email)
  end

  def email=(value)
    @email = value.downcase
  end
end
