class CreateSurveyAnswers < ActiveRecord::Migration[8.0]
  def change
    create_table :survey_answers do |t|
      t.references :survey, null: false, foreign_key: true
      t.references :survey_item, null: true, foreign_key: true
      t.references :user, null: false, foreign_key: true
      t.string :state

      t.timestamps
    end

    add_index :survey_answers, [ :survey_id, :user_id ], unique: true
  end
end
