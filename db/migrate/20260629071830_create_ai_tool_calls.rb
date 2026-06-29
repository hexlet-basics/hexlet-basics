class CreateAiToolCalls < ActiveRecord::Migration[8.1]
  def change
    create_table :ai_tool_calls do |t|
      t.string :tool_call_id, null: false
      t.string :name, null: false
      t.text :thought_signature

      t.jsonb :arguments, default: {}

      t.timestamps
    end

    add_index :ai_tool_calls, :tool_call_id, unique: true
    add_index :ai_tool_calls, :name
  end
end
