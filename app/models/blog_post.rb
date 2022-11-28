# frozen_string_literal: true

class BlogPost < ApplicationRecord
  extend Enumerize

  enumerize :locale, in: I18n.available_locales

  validates :name, presence: true
  validates :slug, presence: true, uniqueness: true
  validates :locale, presence: true
  validates :body, presence: true

  belongs_to :language, optional: true
  belongs_to :creator, class_name: 'User'

  def to_s
    name
  end
end
