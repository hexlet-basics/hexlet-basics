class UserExistValidator
  def validate(record)
    record.errors.add(:email, :user_does_not_exist) unless record.user
  end
end
