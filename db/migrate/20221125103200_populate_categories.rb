class PopulateCategories < ActiveRecord::Migration[7.0]
  def change
    category = Language::Category.create!(slug: :programming, name_en: 'Programming', name_ru: 'Программирование')
    I18n.available_locales.each do |locale|
      Language.with_progress(:completed)
        .joins(current_version: :infos)
        .merge(Language::Version::Info.with_locale(locale))
        .update_all(category_id: category.id)
    end
  end
end
