class AddReferrenceCurrentVersionToLanguages < ActiveRecord::Migration[6.0]
  def change
    add_reference :languages, :current_version, null: true, foreign_key: {to_table: :language_versions}
  end
end
