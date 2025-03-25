class AddFooterToLanguageLandingPages < ActiveRecord::Migration[8.0]
  def change
    add_column :language_landing_pages, :footer, :boolean
    add_column :language_landing_pages, :footer_name, :string
  end
end
