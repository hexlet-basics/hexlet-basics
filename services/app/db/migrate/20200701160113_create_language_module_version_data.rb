class CreateLanguageModuleVersionData < ActiveRecord::Migration[6.0]
  def change
    create_table :language_module_version_infos do |t|
      t.string :name
      t.string :description
      t.string :locale
      t.references :language, null: false, foreign_key: true
      t.references :language_module_version, null: false, foreign_key: true, index: false
      t.references :language_version, null: false, foreign_key: true


      t.timestamps
    end
  end
end
