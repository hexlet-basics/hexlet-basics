class AddCreatorToBlogPost < ActiveRecord::Migration[7.0]
  def change
    add_reference :blog_posts, :creator, null: false, foreign_key: { to_table: :users }
  end
end
