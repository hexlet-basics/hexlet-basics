class CreateLanguageModuleLessonMembers < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_lesson_members do |t|
      t.string :state
      t.references :user, null: false, foreign_key: true
      t.references :language_module_lesson, null: false, foreign_key: true, index: { name: :index_language_module_lesson_member_on_lesson_id }
      t.references :language, null: false, foreign_key: true


      t.timestamps
    end
  end
end
