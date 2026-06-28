# typed: strict

class RemindPasswordForm
  extend T::Sig
  include ActiveFormModel::Virtual

  fields :email

  validates :email, presence: true
  validate :user_exists

  sig { void }
  def user_exists
    errors.add(:email, :user_does_not_exist) if email.present? && !user
  end

  sig { returns(T.nilable(User)) }
  def user
    @user ||= T.let(User.find_by(email: email), T.nilable(User))
  end
end
