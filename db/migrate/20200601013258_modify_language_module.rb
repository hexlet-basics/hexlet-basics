class ModifyLanguageModule < ActiveRecord::Migration[6.0]
  def change
    remove_column :language_modules, :order
    add_reference :language_modules, :current_version, index: true
  end
end
