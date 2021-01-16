# frozen_string_literal: true

class User::ProfileForm < User
  include ActiveFormModel

  permit :first_name, :last_name, :nickname

  # TODO Move to custom validator
  validates :first_name, length: { maximum: 40 },
                         format: { with: UsefulRegexp.without_spec_chars },
                         allow_blank: true

  validates :last_name, length: { maximum: 40 },
                        format: { with: UsefulRegexp.without_spec_chars },
                        allow_blank: true

  validates :nickname, length: { maximum: 40 },
                       format: { with: UsefulRegexp.without_spec_chars },
                       allow_blank: true
end
