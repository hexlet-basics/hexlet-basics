class AddLocaleToReview < ActiveRecord::Migration[7.0]
  def change
    add_column :reviews, :locale, :string
  end
end
