# frozen_string_literal: true

class User::ProfileForm < User
  include ActiveFormModel

  permit :first_name, :last_name, :nickname
end
