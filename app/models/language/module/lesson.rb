class Language::Module::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :module
end
