class ChangeLanguageCurrentVersionColumnNull < ActiveRecord::Migration[6.0]
  def change
    change_column_null :languages, :current_version_id, true
  end
end
