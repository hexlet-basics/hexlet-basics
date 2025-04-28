class AddSurveyItemToSurveys < ActiveRecord::Migration[8.0]
  def change
    add_reference :surveys, :parent_survey_item, null: true, foreign_key: { to_table: :survey_items }
    add_reference :surveys, :parent_survey, null: true, foreign_key: { to_table: :surveys }
  end
end
