# typed: strict

# == Schema Information
#
# Table name: language_lesson_reviews
#
#  id                              :bigint           not null, primary key
#  locale                          :string           not null
#  summary                         :text             default(""), not null
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
  belongs_to :lesson, foreign_key: "language_lesson_id", inverse_of: :reviews
  # у Version и Info нет обратных ассоциаций на ревью
  belongs_to :lesson_version, foreign_key: "language_lesson_version_id", class_name: "Language::Lesson::Version", inverse_of: false
  belongs_to :lesson_version_info, foreign_key: "language_lesson_version_info_id", class_name: "Language::Lesson::Version::Info", inverse_of: false

  # validates :summary, presence: true
  validates :lesson, uniqueness: { scope: :locale }
  validates :locale, presence: true

  # A summary is filled only when a lesson actually has student questions
  # (see ReviewLessonJob), so a non-empty summary marks reviews worth looking at.
  scope :with_summary, -> { where.not(summary: "") }

  sig { params(_auth_object: T.untyped).returns(T.untyped) }
  def self.ransackable_attributes(_auth_object = nil)
    [ "created_at", "id", "language_id", "language_lesson_id", "language_lesson_version_id", "language_lesson_version_info_id", "summary", "updated_at" ]
  end

  sig { params(_auth_object: T.untyped).returns(T.untyped) }
  def self.ransackable_associations(_auth_object = nil)
    [ "language", "lesson" ]
  end
end
