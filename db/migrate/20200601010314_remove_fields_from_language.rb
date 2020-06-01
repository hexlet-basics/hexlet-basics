class RemoveFieldsFromLanguage < ActiveRecord::Migration[6.0]
  def change
    remove_column :languages, :name
    remove_column :languages, :extension
    remove_column :languages, :docker_image
    remove_column :languages, :exercise_filename
    remove_column :languages, :exercise_test_filename
  end
end
