class AddSecondCategory < ActiveRecord::Migration[7.0]
  def change
    Language::Category.create!(slug: :layouting, name_en: 'Layouting', name_ru: 'Верстка')
  end
end
