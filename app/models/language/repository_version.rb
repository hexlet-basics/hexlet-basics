# frozen_string_literal: true

class Language::RepositoryVersion < ApplicationRecord
  belongs_to :language
end
