class AddLearnAsToLanguage < ActiveRecord::Migration[6.1]
  def change
    add_column :languages, :learn_as, :string
  end
end
