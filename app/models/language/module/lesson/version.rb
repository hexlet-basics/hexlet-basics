class Language::Module::Lesson::Version < ApplicationRecord
  belongs_to :language_version
  belongs_to :module_version
end
