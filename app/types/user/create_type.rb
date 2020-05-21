# frozen_string_literal: true

class User::CreateType < User
  include BaseType

  def email=(value)
    super(value.downcase)
  end
end
