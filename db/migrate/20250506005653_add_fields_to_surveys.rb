class AddFieldsToSurveys < ActiveRecord::Migration[8.0]
  def change
    add_column :surveys, :run_always, :boolean, default: false
    add_column :surveys, :run_after_finishing_lessons_count, :integer, default: 0
  end
end
