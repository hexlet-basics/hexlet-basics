class CreateLanguageCategoryQnaItems < ActiveRecord::Migration[8.0]
  def change
    create_table :language_category_qna_items do |t|
      t.references :language_category, null: false, foreign_key: true
      t.string :question
      t.string :answer

      t.timestamps
    end
  end
end
