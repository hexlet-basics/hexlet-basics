# == Schema Information
#
# Table name: surveys
#
#  id                                :bigint           not null, primary key
#  description                       :string
#  locale                            :string
#  question                          :string
#  run_after_finishing_lessons_count :integer          default(0)
#  run_always                        :boolean          default(FALSE)
#  slug                              :string
#  state                             :string
#  created_at                        :datetime         not null
#  updated_at                        :datetime         not null
#  parent_survey_id                  :bigint
#  parent_survey_item_id             :bigint
#
# Indexes
#
#  index_surveys_on_parent_survey_id       (parent_survey_id)
#  index_surveys_on_parent_survey_item_id  (parent_survey_item_id)
#  index_surveys_on_slug_and_locale        (slug,locale) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (parent_survey_id => surveys.id)
#  fk_rails_...  (parent_survey_item_id => survey_items.id)
#
FactoryBot.define do
  factory :survey do
    question { Faker::Lorem.question }
    slug { Faker::Internet.slug }
    # locale { :ru }
    description { Faker::Lorem.paragraph_by_chars }
  end
end
