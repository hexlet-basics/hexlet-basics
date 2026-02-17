class LeadResource < ApplicationResource
  typelize_from Lead

  attributes :id, :user_id, :phone, :whatsapp, :telegram, :survey_answers_data, :courses_data, :created_at

  typelize :string, nullable: true
  attribute :full_name do
    [ it.user.first_name, it.user.last_name ].join " "
  end

  typelize :string, nullable: true
  attribute :email do
    it.user.email
  end

  # typelize "Array<[string, string]>"
  # attribute :state_events do |obj|
  #   obj.aasm.events_for_select
  # end
end
