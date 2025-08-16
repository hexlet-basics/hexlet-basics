# Add your own tasks in files placed in lib/tasks ending in .rake,
# for example lib/tasks/capistrano.rake, and they will automatically be available to Rake.

require_relative "config/application"

Rails.application.load_tasks
# Update js-routes file before javascript build
# task "assets:precompile" => "js:routes"

Rake::Task["db:fixtures:load"].enhance do
  puts "Reset cache counters"

  Language::Member.counter_culture_fix_counts
  Language::Lesson.counter_culture_fix_counts
  # Language::Lesson::Member.counter_culture_fix_counts
  Language::Lesson::Member::Message.counter_culture_fix_counts
  BlogPost::RelatedLanguageItem.counter_culture_fix_counts

  include ActionDispatch::TestProcess

  puts "Upload Images"

  BlogPost.find_each do |p|
    p.cover = File.open(Rails.root.join("test/fixtures/files/blog-post-cover.jpg"))
    p.save!
  end

  Language.find_each do |l|
    l.cover = File.open(Rails.root.join("test/fixtures/files/course-cover.png"))
    l.save!
  end

  puts "Update Models"

  post = BlogPost.find_by!(slug: "full-python-ru") # или другой идентификатор
  post.body = TestHelpers.read_fixture_file("post-about-python.txt")
  post.save!
end
