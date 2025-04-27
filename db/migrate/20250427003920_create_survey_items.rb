class CreateSurveyItems < ActiveRecord::Migration[8.0]
  def change
    create_table :survey_items do |t|
      t.references :survey, null: false, foreign_key: true
      t.integer :order, null: false
      t.string :value
      t.string :state

      t.timestamps
    end
  end
end
