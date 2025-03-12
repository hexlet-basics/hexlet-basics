class AddPinnedToReviews < ActiveRecord::Migration[8.0]
  def change
    add_column :reviews, :pinned, :boolean
  end
end
