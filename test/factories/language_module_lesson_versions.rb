FactoryBot.define do
  factory :language_module_lesson_version do
    language_version { nil }
    language_module_version { nil }
    order { "MyString" }
    original_code { "MyString" }
    prepared_code { "MyString" }
    test_code { "MyString" }
    path_to_code { "MyString" }
  end
end
