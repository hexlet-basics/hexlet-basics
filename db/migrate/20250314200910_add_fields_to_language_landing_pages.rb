class AddFieldsToLanguageLandingPages < ActiveRecord::Migration[8.0]
  def change
    add_column :language_landing_pages, :used_in_header, :string
    add_column :language_landing_pages, :used_in_description, :string
    add_column :language_landing_pages, :outcomes_header, :string
    add_column :language_landing_pages, :outcomes_description, :string
  end
end
