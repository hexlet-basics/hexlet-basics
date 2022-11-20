class CreateBlogPosts < ActiveRecord::Migration[7.0]
  def change
    create_table :blog_posts do |t|
      t.references :course, null: false, foreign_key: true
      t.string :locale
      t.string :state
      t.string :slug
      t.string :name
      t.text :body

      t.timestamps
    end
  end
end
