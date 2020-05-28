class Language::Module::Lesson::Version < ApplicationRecord
  belongs_to :lesson
  belongs_to :language
end
