# frozen_string_literal: true
# typed: true

class SignInForm
  include ActiveFormModel::Virtual

  extend T::Sig

  fields :email, :password

  validates :email, presence: true
  validates :password, presence: true
  validate :authenticate_user

  def authenticate_user
    return if email.blank? || password.blank?

    @user = User.authenticate_by(email: normalized_email, password:)&.then do |authenticated_user|
      authenticated_user if authenticated_user.active?
    end

    return if @user

    errors.add(:password, :cannot_sign_in)
  end

  sig { returns(T.nilable(User)) }
  def user
    @user
  end

  private

    sig { returns(String) }
    def normalized_email
      email.to_s.strip.downcase
    end
end
