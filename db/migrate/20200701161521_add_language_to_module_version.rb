class AddLanguageToModuleVersion < ActiveRecord::Migration[6.0]
  def change
    add_reference :language_module_versions, :language, index: true, foreign_key: true, null: false
  end
end
