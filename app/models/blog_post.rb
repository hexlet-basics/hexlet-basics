class BlogPost < ApplicationRecord
  validates :name, presence: true
  validates :slug, presence: true, uniqueness: true
  validates :locale, presence: true
  validates :body, presence: true

  belongs_to :language, optional: true
end
