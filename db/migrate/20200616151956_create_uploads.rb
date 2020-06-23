class CreateUploads < ActiveRecord::Migration[6.0]
  def change
    create_table :language_uploads do |t|
      t.string :state
      t.string :uploader
      t.string :result
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
