class AddDetailsToModulesAndLessons < ActiveRecord::Migration[6.0]
  def change
    add_column :language_module_lessons, :order, :integer
    add_column :language_modules, :order, :integer
  end
end
