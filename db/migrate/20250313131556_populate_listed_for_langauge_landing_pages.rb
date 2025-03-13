class PopulateListedForLangaugeLandingPages < ActiveRecord::Migration[8.0]
  def change
    Language::LandingPage.find_each do |p|
      p.listed = true
      p.save!
    end
  end
end
