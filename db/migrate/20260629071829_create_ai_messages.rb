class CreateAiMessages < ActiveRecord::Migration[8.1]
  def change
    create_table :ai_messages do |t|
      t.string :role, null: false
      t.text :content
      t.json :content_raw
      t.text :thinking_text
      t.text :thinking_signature
      t.integer :thinking_tokens
      t.integer :input_tokens
      t.integer :output_tokens
      t.integer :cached_tokens
      t.integer :cache_creation_tokens
      t.timestamps
    end

    add_index :ai_messages, :role
  end
end
