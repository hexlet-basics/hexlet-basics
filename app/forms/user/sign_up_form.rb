class User::SignUpForm < User
  include ActiveFormModel

  permit :email, :password, :first_name

  validates :password, presence: true, length: { minimum: 6 }
end
