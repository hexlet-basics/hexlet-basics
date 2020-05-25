class Language::Module::Description < ApplicationRecord
  belongs_to :module
  belongs_to :language
end
