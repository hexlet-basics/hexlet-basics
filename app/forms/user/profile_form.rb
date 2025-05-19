# frozen_string_literal: true

class User::ProfileForm < User
  include ActiveFormModel

  permit :first_name, :last_name, :nickname, :contact_method, :contact_value

  def contact_value=(v)
    write_attribute(:contact_value, v.strip)
  end
end
