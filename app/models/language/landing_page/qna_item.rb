# == Schema Information
#
# Table name: language_landing_page_qna_items
#
#  id                       :integer          not null, primary key
#  answer                   :string
#  question                 :string
#  created_at               :datetime         not null
#  updated_at               :datetime         not null
#  language_landing_page_id :integer          not null
#
# Indexes
#
#  idx_on_language_landing_page_id_98023e1f90  (language_landing_page_id)
#
# Foreign Keys
#
#  language_landing_page_id  (language_landing_page_id => language_landing_pages.id)
#
class Language::LandingPage::QnaItem < ApplicationRecord
  belongs_to :landing_page

  validates :question, presence: true
  validates :answer, presence: true
end
