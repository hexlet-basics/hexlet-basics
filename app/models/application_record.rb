# frozen_string_literal: true

class ApplicationRecord < ActiveRecord::Base
  self.abstract_class = true

  def self.enum_as_hashes(enum_name)
    send(enum_name).map { |key, value| { key: key, value: value } }
  end
end
