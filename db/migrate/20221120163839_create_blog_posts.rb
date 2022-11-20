class CreateBlogPosts < ActiveRecord::Migration[7.0]
  def change
    create_table :blog_posts do |t|
      t.references :language, null: true, foreign_key: true
      t.string :locale
      t.string :state
      t.string :slug
      t.string :name
      t.text :body

      t.index :slug, unique: true

      t.timestamps
    end
  end
end
