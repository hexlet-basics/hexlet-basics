# frozen_string_literal: true

# == Schema Information
#
# Table name: languages
#
#  id                     :bigint           not null, primary key
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
#  fk_rails_...              (category_id => language_categories.id)
#  fk_rails_...              (current_version_id => language_versions.id)
#  languages_upload_id_fkey  (upload_id => uploads.id)
#
class Language < ApplicationRecord
  include LanguageRepository
  include LanguagePresenter

  extend Enumerize

  def self.ransackable_attributes(_auth_object = nil)
    ['created_at']
  end

  enumerize :slug, in: %i[ada bash clang clojure cobol cpp csharp css crystal dart dlang elixir elm fortran go groovy haskell html java
                          javascript kotlin lua objectivec ocaml perl php prolog python racket rescript rproject ruby
                          rust scala smalltalk swift typescript perl powershell ocaml layout-designer pre-course-java pre-course-python
                          pre-course-javascript]

  enumerize :progress, in: %i[completed in_development draft], default: :draft, scope: true, predicates: { prefix: true }
  # TODO: move to language version and populate inside the job
  enumerize :learn_as, in: %i[first_language second_language], default: :first_language

  validates :slug, presence: true
  # validates :learn_as, presence: true

  belongs_to :current_version, optional: true, class_name: 'Language::Version'
  belongs_to :category, optional: true

  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :versions, dependent: :destroy, class_name: '::Language::Version'
  has_many :infos, class_name: '::Language::Version::Info', dependent: :restrict_with_exception
  has_many :members, dependent: :destroy
  has_many :reviews, dependent: :restrict_with_exception
  has_many :blog_posts, dependent: :restrict_with_exception

  has_many :current_module_infos, through: :current_version, source: :module_infos
  has_many :current_lesson_infos, through: :current_version, source: :lesson_infos
  has_many :current_lesson_versions, through: :current_version, source: :lesson_versions
  has_many :current_module_versions, through: :current_version, source: :module_versions
  has_many :current_lessons, through: :current_version, source: :lessons

  delegate :to_s, to: :current_version
  delegate :to_hash, to: :current_version
  delegate :name, to: :current_version

  def duration
    # TODO Пересадить на counter_culture от Language::Version
    current_lessons.size * 15 / 60
  end
end
