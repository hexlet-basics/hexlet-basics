FactoryBot.define do
  factory :blog_post_related_language_item, class: "BlogPost::RelatedLanguageItem" do
    blog_post { nil }
    language { nil }
  end
end
