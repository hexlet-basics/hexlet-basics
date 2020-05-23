class Language::Module::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :language_module, class_name: 'Language::Module'

  validates :slug, uniqueness: { scope: [:language, :language_module], message: 'slug should be uniqueness on language and module' },
    presence: true
end
