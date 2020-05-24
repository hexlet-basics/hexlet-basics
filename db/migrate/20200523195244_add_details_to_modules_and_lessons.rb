class AddDetailsToModulesAndLessons < ActiveRecord::Migration[6.0]
  def change
    add_column :language_module_lessons, :order, :integer
    add_column :language_modules, :order, :integer

    add_column :languages, :docker_image, :string
    add_column :languages, :exercise_filename, :string
    add_column :languages, :exercise_test_filename, :string
  end
end
