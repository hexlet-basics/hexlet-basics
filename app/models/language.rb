# frozen_string_literal: true

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

  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :versions, dependent: :destroy, class_name: '::Language::Version'
  has_many :infos, class_name: '::Language::Version::Info', dependent: :restrict_with_exception
  has_many :members, dependent: :destroy

  has_many :current_module_infos, through: :current_version, source: :module_infos
  has_many :current_lesson_infos, through: :current_version, source: :lesson_infos
  has_many :current_lesson_versions, through: :current_version, source: :lesson_versions
  has_many :current_module_versions, through: :current_version, source: :module_versions
  has_many :current_lessons, through: :current_version, source: :lessons

  delegate :to_s, to: :current_version
  delegate :to_hash, to: :current_version
  delegate :name, to: :current_version

  def duration
    lessons_count * 15 / 60
  end
end
