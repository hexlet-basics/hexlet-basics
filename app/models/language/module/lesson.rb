# frozen_string_literal: true

class Language::Module::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :language_module, class_name: 'Language::Module'

  validates :slug, uniqueness: { scope: %i[language language_module], message: 'slug should be uniqueness on language and module' },
                   presence: true
  validates :path_to_code, presence: true
  validates :order, presence: true
end
