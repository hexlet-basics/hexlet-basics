class CreateLanguageRepositoryVersions < ActiveRecord::Migration[6.0]
  def change
    create_table :language_repository_versions do |t|
      t.string :language_name
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
