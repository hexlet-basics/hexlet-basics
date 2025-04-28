class CreateUserSurveyPivots < ActiveRecord::Migration[8.0]
  def change
    create_table :user_survey_pivots do |t|
      t.references :user, null: false, foreign_key: true
      t.references :coding_experience_item, null: true, foreign_key: { to_table: :survey_items }
      t.references :goal_item, null: true, foreign_key: { to_table: :survey_items }

      t.timestamps
    end
  end
end
