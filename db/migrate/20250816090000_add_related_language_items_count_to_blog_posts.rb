class AddRelatedLanguageItemsCountToBlogPosts < ActiveRecord::Migration[8.0]
  def up
    add_column :blog_posts, :related_language_items_count, :integer, null: false, default: 0

    execute <<~SQL
      UPDATE blog_posts
      SET related_language_items_count = sub.count
      FROM (
        SELECT blog_post_id, COUNT(*) AS count
        FROM blog_post_related_language_items
        GROUP BY blog_post_id
      ) AS sub
      WHERE blog_posts.id = sub.blog_post_id
    SQL
  end

  def down
    remove_column :blog_posts, :related_language_items_count
  end
end
