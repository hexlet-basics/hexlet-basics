# typed: false
# frozen_string_literal: true

# == Schema Information
#
# Table name: banners
#
#  id          :bigint           not null, primary key
#  background  :string           default("cta_gradient"), not null
#  body        :text             not null
#  finishes_at :datetime
#  locale      :string           not null
#  starts_at   :datetime
#  state       :string           default("draft"), not null
#  url         :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
# Indexes
#
#  index_banners_on_locale_and_state  (locale,state)
#
FactoryBot.define do
  factory :banner do
    locale { "ru" }
    body { "Курс **со скидкой 17%** • только до конца месяца" }
    url { "https://example.com/promo" }
    background { "cta_gradient" }
    state { "published" }
    starts_at { 1.day.ago }
    finishes_at { 1.day.from_now }
  end
end
