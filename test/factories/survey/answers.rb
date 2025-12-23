# == Schema Information
#
# Table name: survey_answers
#
#  id             :bigint           not null, primary key
#  state          :string
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  survey_id      :bigint           not null
#  survey_item_id :bigint
#  user_id        :bigint           not null
#
# Indexes
#
#  index_survey_answers_on_survey_id              (survey_id)
#  index_survey_answers_on_survey_id_and_user_id  (survey_id,user_id) UNIQUE
#  index_survey_answers_on_survey_item_id         (survey_item_id)
#  index_survey_answers_on_user_id                (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (survey_id => surveys.id)
#  fk_rails_...  (survey_item_id => survey_items.id)
#  fk_rails_...  (user_id => users.id)
#
FactoryBot.define do
  factory :survey_answer, class: "Survey::Answer" do
    survey { nil }
    survey_answer { nil }
    user { nil }
  end
end
