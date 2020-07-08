class CreateLanguageModuleVersions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_versions do |t|
      t.references :language, null: false, foreign_key: true
      t.references :language_version, null: false, foreign_key: true
      t.references :language_module, null: false, foreign_key: true
      t.string :order

      t.timestamps
    end
  end
end
