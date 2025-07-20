class CreateBlogPostRelatedLanguageItems < ActiveRecord::Migration[8.0]
  def change
    create_table :blog_post_related_language_items do |t|
      t.references :blog_post, null: false, foreign_key: true
      t.references :language, null: false, foreign_key: true

      t.timestamps
    end
  end
end
