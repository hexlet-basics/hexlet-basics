# == Schema Information
#
# Table name: language_category_qna_items
#
#  id                   :bigint           not null, primary key
#  answer               :string
#  question             :string
#  created_at           :datetime         not null
#  updated_at           :datetime         not null
#  language_category_id :bigint           not null
#
# Indexes
#
#  index_language_category_qna_items_on_language_category_id  (language_category_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_category_id => language_categories.id)
#
class Language::Category::QnaItem < ApplicationRecord
  belongs_to :language_category, class_name: "Language::Category"

  validates :question, presence: true
  validates :answer, presence: true
end
