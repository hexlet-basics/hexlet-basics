class CreateLanguageVersionInfos < ActiveRecord::Migration[6.1]
  def change
    create_table :language_version_infos do |t|
      t.references :language, null: false, foreign_key: true
      t.references :language_version, null: false, foreign_key: true
      t.string :locale
      t.string :description

      t.timestamps
    end
  end
end
