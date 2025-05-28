class AddHexletProgramUrlToLanguages < ActiveRecord::Migration[8.0]
  def change
    add_column :languages, :hexlet_program_landing_page, :string, null: true
  end
end
