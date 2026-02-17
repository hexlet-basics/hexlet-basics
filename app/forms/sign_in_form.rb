# frozen_string_literal: true
# typed: true

class SignInForm
  include ActiveFormModel::Virtual

  extend T::Sig

  fields :email, :password

  validates :email, presence: true
  validates :password, presence: true
  validate :user_exists, :user_can_sign_in

  def user_can_sign_in
    errors.add(:password, :cannot_sign_in) if password.present? && user && !user.valid_password?(password)
  end

  def user_exists
    errors.add(:email, :user_does_not_exist_html) unless user
  end

  sig { returns(T.nilable(User)) }
  def user
    @user ||= User.active.find_by(email: email) if email.present?
  end
end
