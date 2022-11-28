class AddTitleToBlogPost < ActiveRecord::Migration[7.0]
  def change
    add_column :blog_posts, :title, :string
  end
end
