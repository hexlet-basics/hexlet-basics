class CreateLanguages < ActiveRecord::Migration[6.0]
  def change
    create_table :languages do |t|
      t.string :name
      t.string :slug
      t.string :extension

      t.timestamps
    end
  end
end
