class CreateSurveys < ActiveRecord::Migration[8.0]
  def change
    create_table :surveys do |t|
      t.string :question
      t.string :state
      t.string :slug
      t.string :locale
      t.string :description

      t.timestamps
    end

    add_index :surveys, [ :slug, :locale ], unique: true
  end
end
