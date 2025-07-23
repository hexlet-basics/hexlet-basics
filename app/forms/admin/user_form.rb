class Admin::UserForm < User
  include ActiveFormModel

  permit :email, :first_name, :last_name, :admin
end
