# == Schema Information
#
# Table name: language_lesson_member_messages
#
#  id                        :bigint           not null, primary key
#  body                      :text
#  role                      :string
#  created_at                :datetime         not null
#  updated_at                :datetime         not null
#  language_id               :integer          not null
#  language_lesson_id        :integer          not null
#  language_lesson_member_id :integer          not null
#  user_id                   :bigint
#
# Indexes
#
#  idx_on_language_lesson_member_id_fe254654e9                  (language_lesson_member_id)
#  index_language_lesson_member_messages_on_language_id         (language_id)
#  index_language_lesson_member_messages_on_language_lesson_id  (language_lesson_id)
#  index_language_lesson_member_messages_on_user_id             (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (language_lesson_id => language_lessons.id)
#  fk_rails_...  (language_lesson_member_id => language_lesson_members.id)
#  fk_rails_...  (user_id => users.id)
#
require "test_helper"

class Language::Lesson::Member::MessageTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
