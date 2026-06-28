# typed: false
# frozen_string_literal: true

class CreateBanners < ActiveRecord::Migration[8.0]
  def change
    create_table :banners do |t|
      t.string :locale, null: false
      t.text :body, null: false
      t.string :url
      t.string :background, null: false, default: "cta_gradient"
      t.string :state, null: false, default: "draft"
      t.datetime :starts_at
      t.datetime :finishes_at

      t.timestamps
    end

    add_index :banners, %i[locale state]
  end
end
