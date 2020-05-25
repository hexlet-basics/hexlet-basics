class Language::Module::Lesson::Description < ApplicationRecord
  belongs_to :lesson
  belongs_to :language
end
