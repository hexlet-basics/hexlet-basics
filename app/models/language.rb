# frozen_string_literal: true

# == Schema Information
#
# Table name: languages
#
#  id                     :integer          not null, primary key
#  docker_image           :string(255)
#  exercise_filename      :string(255)
#  exercise_test_filename :string(255)
#  extension              :string(255)
#  learn_as               :string
#  lessons_count          :integer          default(0), not null
#  members_count          :integer          default(0), not null
#  name                   :string(255)
#  order                  :integer
#  progress               :string
#  slug                   :string(255)
#  state                  :string(255)
#  created_at             :datetime         not null
#  updated_at             :datetime         not null
#  category_id            :bigint
#  current_version_id     :bigint
#  upload_id              :bigint
#
# Indexes
#
#  index_languages_on_category_id         (category_id)
#  index_languages_on_current_version_id  (current_version_id)
#  languages_slug_index                   (slug) UNIQUE
#  languages_upload_id_index              (upload_id)
#
# Foreign Keys
#
#  category_id         (category_id => language_categories.id)
#  current_version_id  (current_version_id => language_versions.id)
#  upload_id           (upload_id => uploads.id)
#
class Language < ApplicationRecord
  # TODO: remove name
  include LanguageRepository

  has_one_attached :cover do |attachable|
    attachable.variant :list, resize_to_limit: [ 390, 320 ], preprocessed: true
    attachable.variant :thumb, resize_to_limit: [ 39, 32 ], preprocessed: true
  end

  # extend Enumerize

  def self.ransackable_associations(_)
    []
  end

  def self.ransackable_attributes(_auth_object = nil)
    [ "created_at", "learn_as", "progress" ]
  end

  enum :progress, { completed: "completed", in_development: "in_development", draft: "draft" }, default: "draft", suffix: true, validate: true
  enum :learn_as, { first_language: "first_language", second_language: "second_language" }, default: "first_language", suffix: true, validate: true
  # enumerize :progress, in: %i[completed in_development draft], default: :draft, scope: true, predicates: { prefix: true }
  # enumerize :learn_as, in: %i[first_language second_language], default: :first_language

  # NOTE: must be part of the docker image: hexlet-basics/exercises-<slug>
  validates :slug, presence: true
  # validates :learn_as, presence: true

  belongs_to :current_version, optional: true, class_name: "Language::Version"
  # belongs_to :category, optional: true

  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :landing_pages, dependent: :destroy
  has_many :versions, dependent: :destroy, class_name: "::Language::Version"
  # has_many :infos, class_name: "::Language::Version::Info", dependent: :restrict_with_exception
  has_many :members, dependent: :destroy
  has_many :reviews, dependent: :restrict_with_exception
  has_many :blog_posts, dependent: :restrict_with_exception

  has_many :current_module_infos, through: :current_version, source: :module_infos
  has_many :current_lesson_infos, through: :current_version, source: :lesson_infos
  has_many :current_lesson_versions, through: :current_version, source: :lesson_versions
  has_many :current_module_versions, through: :current_version, source: :module_versions
  has_many :current_lessons, through: :current_version, source: :lessons

  # delegate :to_s, to: :current_version
  # delegate :serializable_data, to: :current_version
  # delegate :name, to: :current_version, allow_nil: true

  def duration
    # TODO Пересадить на counter_culture от Language::Version
    lessons_count * 15 / 60
  end

  def repository_url
    "https://github.com/hexlet-basics/exercises-#{slug}"
  end
end
