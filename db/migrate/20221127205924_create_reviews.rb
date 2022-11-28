class CreateReviews < ActiveRecord::Migration[7.0]
  def change
    create_table :reviews do |t|
      t.references :language, null: false, foreign_key: true
      t.references :user, null: false, foreign_key: true
      t.string :state
      t.text :body

      t.timestamps
    end
  end
end
