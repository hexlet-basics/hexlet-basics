class AddOrderToBlogPostRelatedLanguageItems < ActiveRecord::Migration[8.0]
  def change
    add_column :blog_post_related_language_items, :order, :integer, null: true
  end
end
