class CreateLanguageMembers < ActiveRecord::Migration[6.0]
  def change
    create_table :language_members do |t|
      t.references :language, null: false, foreign_key: true
      t.references :user, null: false, foreign_key: true
      t.string :state

      t.timestamps
    end
  end
end
