# frozen_string_literal: true

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
#  language_module_lessons_language_id_fkey  (language_id => languages.id)
#  language_module_lessons_module_id_fkey    (module_id => language_modules.id)
#  language_module_lessons_upload_id_fkey    (upload_id => uploads.id)
#
class Language::Lesson < ApplicationRecord
  include Language::LessonRepository
  include AASM

  # FIXME: add unique index
  validates :slug, presence: true, uniqueness: { scope: :language }

  counter_culture :language

  belongs_to :language
  belongs_to :module

  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :infos, through: :versions, class_name: 'Language::Lesson::Version::Info'

  aasm :state do
    state :created, initial: true
    state :active
    state :archived

    event :activate do
      transitions from: %i[created archived], to: :active
    end

    event :mark_as_archived do
      transitions from: %i[created archived], to: :archived
    end
  end

  def to_s
    slug
  end

  # TODO: dont change defaults
  def to_hash(*_args)
    attrs = attributes.extract! 'id', 'slug', 'created_at'
    attrs.to_hash
  end
end
