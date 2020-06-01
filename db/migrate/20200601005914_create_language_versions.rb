class CreateLanguageVersions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_versions do |t|
      t.string :docker_image
      t.string :exercise_filename
      t.string :exercise_test_filename
      t.string :extension
      t.string :name
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
