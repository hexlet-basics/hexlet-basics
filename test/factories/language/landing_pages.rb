# == Schema Information
#
# Table name: language_landing_pages
#
#  id                          :bigint           not null, primary key
#  description                 :string
#  footer                      :boolean
#  footer_name                 :string
#  header                      :string
#  listed                      :boolean
#  locale                      :string
#  main                        :boolean
#  meta_description            :string
#  meta_title                  :string
#  name                        :string
#  order                       :string
#  outcomes_description        :string
#  outcomes_header             :string
#  slug                        :string
#  state                       :string
#  used_in_description         :string
#  used_in_header              :string
#  created_at                  :datetime         not null
#  updated_at                  :datetime         not null
#  landing_page_to_redirect_id :bigint
#  language_category_id        :integer
#  language_id                 :integer          not null
#
# Indexes
#
#  index_language_landing_pages_on_landing_page_to_redirect_id  (landing_page_to_redirect_id)
#  index_language_landing_pages_on_language_category_id         (language_category_id)
#  index_language_landing_pages_on_language_id                  (language_id)
#
# Foreign Keys
#
#  fk_rails_...  (landing_page_to_redirect_id => language_landing_pages.id)
#  fk_rails_...  (language_category_id => language_categories.id)
#  fk_rails_...  (language_id => languages.id)
#
FactoryBot.define do
  factory :language_landing_page, class: "Language::LandingPage" do
    language { nil }
    title { "MyString" }
    header { "MyString" }
    description { "MyString" }
  end
end
