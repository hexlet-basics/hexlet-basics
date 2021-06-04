# frozen_string_literal: true

FactoryBot.define do
  factory :language_version_info, class: 'Language::Version::Info' do
    language { nil }
    language_version { nil }
    locale { 'MyString' }
    description { 'MyString' }
  end
end
