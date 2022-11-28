class RemoveTitleFromBlogPosts < ActiveRecord::Migration[7.0]
  def change
    remove_column :blog_posts, :title
  end
end
