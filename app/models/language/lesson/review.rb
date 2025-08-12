# == Schema Information
#
# Table name: language_lesson_reviews
#
#  id                              :bigint           not null, primary key
#  summary                         :text             not null
#  created_at                      :datetime         not null
#  updated_at                      :datetime         not null
#  language_id                     :bigint           not null
#  language_lesson_id              :bigint           not null
#  language_lesson_version_id      :bigint           not null
#  language_lesson_version_info_id :bigint           not null
#
# Indexes
#
#  idx_on_language_lesson_version_info_id_e5ef52eeca            (language_lesson_version_info_id)
#  index_language_lesson_reviews_on_language_id                 (language_id)
#  index_language_lesson_reviews_on_language_lesson_id          (language_lesson_id)
#  index_language_lesson_reviews_on_language_lesson_version_id  (language_lesson_version_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (language_lesson_id => language_lessons.id)
#  fk_rails_...  (language_lesson_version_id => language_lesson_versions.id)
#  fk_rails_...  (language_lesson_version_info_id => language_lesson_version_infos.id)
#
class Language::Lesson::Review < ApplicationRecord
  belongs_to :language
  belongs_to :lesson
  belongs_to :lesson_version
  belongs_to :lesson_info

  validates :summary, presence: true

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "id", "language_id", "language_lesson_id", "language_lesson_version_id", "language_lesson_version_info_id", "summary", "updated_at" ]
  end
end
