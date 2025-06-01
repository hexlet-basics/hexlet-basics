# == Schema Information
#
# Table name: language_category_items
#
#  id                       :bigint           not null, primary key
#  created_at               :datetime         not null
#  updated_at               :datetime         not null
#  language_category_id     :bigint           not null
#  language_landing_page_id :bigint           not null
#
# Indexes
#
#  index_language_category_items_on_language_category_id      (language_category_id)
#  index_language_category_items_on_language_landing_page_id  (language_landing_page_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_category_id => language_categories.id)
#  fk_rails_...  (language_landing_page_id => language_landing_pages.id)
#
class Language::Category::Item < ApplicationRecord
  validates :language_laindg_page, uniqueness: { scope: :language_category }

  belongs_to :language_category, class_name: "Language::Category"
  belongs_to :language_landing_page, class_name: "Language::LandingPage"
end
