class AddStateToLanguageLessons < ActiveRecord::Migration[6.0]
  def change
    Language::Lesson.update_all(state: :active)
  end
end
