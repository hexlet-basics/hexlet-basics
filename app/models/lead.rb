# == Schema Information
#
# Table name: leads
#
#  id                  :bigint           not null, primary key
#  email               :string
#  phone               :string
#  state               :string
#  survey_answers_data :text
#  telegram            :string
#  whatsapp            :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  user_id             :bigint           not null
#
# Indexes
#
#  index_leads_on_user_id  (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (user_id => users.id)
#
class Lead < ApplicationRecord
  belongs_to :user
  serialize :survey_answers_data

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "email", "id", "phone", "state", "survey_answers_data", "telegram", "updated_at", "user_id", "whatsapp" ]
  end

  def self.create_or_update(user)
    lead = find_or_initialize_by user: user
    new_lead = lead.new_record?
    if user.contact_method? && user.contact_value?
      lead[user.contact_method] = user.contact_value
    end

    survey_answers_data = []
    user.survey_answers.each do |answer|
      answer_data = {
        question: answer.survey.question,
        answer: answer.survey_item.value
      }

      survey_answers_data << answer_data
    end

    lead.survey_answers_data = survey_answers_data

    lead.save!

    event = nil
    if new_lead
      event_data = lead.slice(
        :user_id,
        :phone,
        :telegram,
        :whatsapp,
        :survey_answers_data
      )

      event = LeadCreatedEvent.new(data: event_data)
    end

    { lead:, event: }
  end
end
