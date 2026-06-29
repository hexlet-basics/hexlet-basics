class CreateAiChats < ActiveRecord::Migration[8.1]
  def change
    create_table :ai_chats do |t|
      t.timestamps
    end
  end
end
