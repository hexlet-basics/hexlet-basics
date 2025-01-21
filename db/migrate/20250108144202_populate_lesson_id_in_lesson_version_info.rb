class PopulateLessonIdInLessonVersionInfo < ActiveRecord::Migration[8.0]
  def change
    Language::Lesson::Version::Info.find_each do |i|
      i.language_lesson = i.version.lesson
      i.save!
    end
  end
end
