class CreateBookRequests < ActiveRecord::Migration[8.0]
  def change
    create_table :book_requests do |t|
      t.references :user, null: false, foreign_key: true, index: { unique: true }

      t.timestamps
    end
  end
end
