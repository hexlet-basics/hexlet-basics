# frozen_string_literal: true

class User::ProfileForm < User
  include ActiveFormModel

  permit :first_name, :last_name

  # TODO Add sefullRegexp.alphanumeric_for_all_langs
  validates :first_name, length: { maximum: 40 }
  validates :last_name, length: { maximum: 40 }
end
