class AddSlugToSurveyItems < ActiveRecord::Migration[8.0]
  def change
    add_column :survey_items, :slug, :string
  end
end
