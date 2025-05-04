class CreateLeads < ActiveRecord::Migration[8.0]
  def change
    create_table :leads do |t|
      t.references :user, null: false, foreign_key: true
      t.string :state
      t.string :email
      t.string :phone
      t.string :telegram
      t.string :whatsapp
      t.text :survey_answers_data

      t.timestamps
    end
  end
end
