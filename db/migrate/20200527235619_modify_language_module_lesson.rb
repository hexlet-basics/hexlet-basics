class ModifyLanguageModuleLesson < ActiveRecord::Migration[6.0]
  def change
    remove_column :language_module_lessons, :original_code
    remove_column :language_module_lessons, :prepared_code
    remove_column :language_module_lessons, :test_code
    remove_column :language_module_lessons, :path_to_code
    remove_column :language_module_lessons, :natural_order
    add_column :language_module_lessons, :state, :string
  end
end
