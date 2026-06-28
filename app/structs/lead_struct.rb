# typed: strict
# frozen_string_literal: true

class LeadStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :contact_method, T.nilable(String), default: nil
  const :contact_value, T.nilable(String), default: nil
  const :ym_client_id, T.nilable(String), default: nil

  validates :contact_method, presence: true
  validates :contact_value, presence: true
end
