class AddFieldsToLanguage < ActiveRecord::Migration[6.0]
  def change
    add_column :languages, :docker_image, :string
    add_column :languages, :exercise_filename, :string
    add_column :languages, :exercise_test_filename, :string
  end
end
