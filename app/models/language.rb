# frozen_string_literal: true

# == Schema Information
#
# Table name: languages
#
#  id                     :bigint           not null, primary key
#  name                   :string(255)
#  slug                   :string(255)
#  extension              :string(255)
#  docker_image           :string(255)
#  exercise_filename      :string(255)
#  exercise_test_filename :string(255)
#  state                  :string(255)
#  upload_id              :bigint
#  updated_at             :datetime         not null
#  progress               :string
#  current_version_id     :bigint
#  created_at             :datetime         not null
#  learn_as               :string
#  lessons_count          :integer          default(0), not null
#  members_count          :integer          default(0), not null
#  order                  :integer
#  category_id            :bigint
#
class Language < ApplicationRecord
  include LanguageRepository
  include LanguagePresenter

  extend Enumerize

  enumerize :slug, in: %i[ada bash clang clojure cobol cpp csharp css dart elixir elm fortran go groovy haskell html java
                          javascript kotlin lua objectivec ocaml perl php prolog python racket rescript rproject ruby
                          rust scala smalltalk swift typescript]

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
  delegate :locales, to: :current_version

  def duration
    lessons_count * 15 / 60
  end
end
