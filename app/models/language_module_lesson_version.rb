class LanguageModuleLessonVersion < ApplicationRecord
  belongs_to :language_version
  belongs_to :language_module_version
end
