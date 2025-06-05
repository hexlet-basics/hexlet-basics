class LandingPage < ActiveRecord::Migration[8.0]
  def change
    # t.references :scenario, null: false, foreign_key: { to_table: :survey_scenarios }
    add_reference :language_landing_pages, :landing_page_to_redirect, foreign_key: { to_table: :language_landing_pages }
  end
end
