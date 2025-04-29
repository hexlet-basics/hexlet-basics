class AddAssistantMessagesCountToUsers < ActiveRecord::Migration[8.0]
  def change
    add_column :users, :assistant_messages_count, :integer, default: 0
  end
end
