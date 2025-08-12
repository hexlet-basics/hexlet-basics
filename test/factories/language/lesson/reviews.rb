FactoryBot.define do
  factory :language_lesson_review, class: "Language::Lesson::Review" do
    language { nil }
    lesson { nil }
    lesson_version { nil }
    lesson_info { nil }
    summary { "MyText" }
  end
end
