class PopulateLandingPages < ActiveRecord::Migration[8.0]
  def up
    Language::Version::Info.find_each do |i|
      p = Language::LandingPage.new
      p.language = i.language
      p.locale = i.locale
      p.slug = i.language.slug
      p.main = true
      p.language_category = i.language.category
      p.header = i.header
      p.meta_title = i.title
      p.meta_description = i.seo_description
      p.description = i.description
      p.save
    end
  end
end
