class Language::Lesson::Review < ApplicationRecord
  belongs_to :language
  belongs_to :lesson
  belongs_to :lesson_version
  belongs_to :lesson_info

  validates :summary, precense: true
end
