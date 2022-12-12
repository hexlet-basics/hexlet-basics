# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lessons
#
#  id            :bigint           not null, primary key
#  slug          :string(255)
#  state         :string(255)
#  order         :integer
#  original_code :text
#  prepared_code :text
#  test_code     :text
#  path_to_code  :string(255)
#  module_id     :bigint
#  language_id   :bigint
#  upload_id     :bigint
#  updated_at    :datetime         not null
#  natural_order :integer
#  created_at    :datetime         not null
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
