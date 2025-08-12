# == Schema Information
#
# Table name: language_lessons
#
#  id            :bigint           not null, primary key
#  natural_order :integer
#  order         :integer
#  original_code :text
#  path_to_code  :string(255)
#  prepared_code :text
#  slug          :string(255)
#  state         :string(255)
#  test_code     :text
#  created_at    :datetime         not null
#  updated_at    :datetime         not null
#  language_id   :bigint
#  module_id     :bigint
#  upload_id     :bigint
#
# Indexes
#
#  index_language_lessons_on_language_id_and_slug  (language_id,slug) UNIQUE
#  language_module_lessons_language_id_index       (language_id)
#  language_module_lessons_module_id_index         (module_id)
#  language_module_lessons_upload_id_index         (upload_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (module_id => language_modules.id)
#  fk_rails_...  (upload_id => uploads.id)
#
class Language::Lesson < ApplicationRecord
  include AASM
  include Language::LessonRepository


  # FIXME: add unique index
  validates :slug, presence: true, uniqueness: { scope: :language }

  counter_culture :language

  belongs_to :language
  belongs_to :module
  belongs_to :review, foreign_key: "language_lesson_id", dependent: :destroy, class_name: "Language::Lesson::Review"

  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy
  has_many :messages, through: :members, dependent: :destroy

  has_many :infos, through: :versions, class_name: "Language::Lesson::Version::Info"
  # has_one :localed_info, -> { merge(Language::Lesson::Version::Info.with_locale) },
  #   through: :versions, class_name: "Language::Lesson::Version::Info"
  # has_one :localed_info, -> { merge(Language::Lesson::Version::Info.with_locale) },
  #   class_name: "Language::Lesson::Version::Info", through: :versions, source: :infos

  enum :state, { created: "created", active: "active", archived: "archived" }, suffix: true
  aasm :state, enum: true do
   state :created, initial: true
   state :active
   state :archived

   event :activate do
     transitions from: %i[created archived], to: :active
   end

   event :mark_as_archived do
     transitions from: %i[created active], to: :archived
   end
 end

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "id", "language_id", "module_id", "natural_order", "order", "original_code", "path_to_code", "prepared_code", "review", "slug", "state", "test_code", "updated_at", "upload_id" ]
  end

  def to_s
    slug
  end

  # work when joined using with_localed_info
  def localed_info
    infos.first
  end
end
