class PopulateLanguageCategories < ActiveRecord::Migration[8.0]
  def change
    Language::Category.find_each do |category|
      category.name = category.name_ru
      category.locale = 'ru'
      category.save!
    end
  end
end
