class AddDescriptionToBlogPost < ActiveRecord::Migration[7.0]
  def change
    add_column :blog_posts, :description, :string
  end
end
