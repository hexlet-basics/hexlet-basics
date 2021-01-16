class MigrateLessonMembersLanguage < ActiveRecord::Migration[6.0]
  def change
    languages = Language.all
    languages.each do |language|
      puts language.id.inspect
      lessons = language.lessons

      Language::Lesson::Member.where(lesson: lessons).update_all(language_id: language.id)
    end

    change_column_null :language_lesson_members, :language_id, false

  end
end
