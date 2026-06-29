class DropLegacyOpenaiAssistant < ActiveRecord::Migration[8.1]
  def up
    drop_table :language_lesson_member_messages
    remove_column :languages, :openai_assistant_id
    remove_column :language_lesson_members, :openai_thread_id
  end

  def down
    add_column :language_lesson_members, :openai_thread_id, :string
    add_column :languages, :openai_assistant_id, :string

    create_table :language_lesson_member_messages do |t|
      t.text :body
      t.string :role
      t.references :language_lesson_member, null: false, foreign_key: true
      t.references :language, null: false, foreign_key: true
      t.references :language_lesson, null: false, foreign_key: true
      t.references :user, foreign_key: true
      t.timestamps
    end
  end
end
