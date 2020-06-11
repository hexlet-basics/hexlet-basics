# frozen_string_literal: true

class User::SignUpType < User
  include BaseType

  validates :password, presence: true, length: { minimum: 6 }

  def email=(value)
    super(value.downcase)
  end
end
