class CreateBlogPostLikes < ActiveRecord::Migration[8.0]
  def change
    create_table :blog_post_likes do |t|
      t.references :blog_post, null: false, foreign_key: true
      t.references :user, null: true, foreign_key: true

      t.timestamps
    end
  end
end
