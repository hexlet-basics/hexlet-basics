# frozen_string_literal: true

class AddOrderToLanguages < ActiveRecord::Migration[7.0]
  def change
    add_column :languages, :order, :integer
  end
end
