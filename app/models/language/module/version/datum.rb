# frozen_string_literal: true

class Language::Module::Version::Datum < ApplicationRecord
  belongs_to :language
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :version

  def to_s
    name
  end
end
