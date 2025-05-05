class LeadResource < ApplicationResource
  typelize_from Lead

  attributes :id, :user_id, :phone, :whatsapp, :telegram, :survey_answers_data

  typelize :string, nullable: true
  attribute :full_name do |obj|
    [ obj.user.first_name, obj.user.last_name ].join " "
  end

  typelize :string, nullable: true
  attribute :email do |obj|
    obj.user.email
  end
  #
  # typelize "Array<[string, string]>"
  # attribute :state_events do |obj|
  #   obj.aasm.events_for_select
  # end
end
