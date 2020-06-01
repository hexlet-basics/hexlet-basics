class AddCurrentVersionToLanguage < ActiveRecord::Migration[6.0]
  def change
    add_reference :languages, :current_version, index: true
  end
end
