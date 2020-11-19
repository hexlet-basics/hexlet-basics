# frozen_string_literal: true

class Language::Module::Version::Info < ApplicationRecord
  include Language::Module::Version::InfoRepository

  belongs_to :language
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :version

  def to_s
    name
  end
end
