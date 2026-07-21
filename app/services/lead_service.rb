# typed: strict
# frozen_string_literal: true

class LeadService < ApplicationService
  class << self
    extend T::Sig

    sig { params(struct: LeadStruct, user: User, courses_data: T::Array[T::Hash[Symbol, T.untyped]], ahoy_visit: T.nilable(Ahoy::Visit)).returns(Typed::Result[Lead, Lead]) }
    def create(struct, user:, courses_data:, ahoy_visit:)
      lead = Lead.new(
        ym_client_id: struct.ym_client_id,
        user:,
        courses_data:,
        survey_answers_data: [],
        ahoy_visit:
      )

      # contact_method names the column (telegram/phone/whatsapp) to fill
      lead.write_attribute(struct.contact_method, struct.contact_value) if struct.contact_method.present?

      return fail_with(lead) unless lead.save

      success_with(lead)
    end
  end
end
