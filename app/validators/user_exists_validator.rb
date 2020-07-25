class UserExistsValidator < ActiveModel::Validator
  def validate(record)
    record.errors.add(:email, :user_does_not_exist) if record.user.nil?
  end
end
