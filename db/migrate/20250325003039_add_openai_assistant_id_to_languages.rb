class AddOpenaiAssistantIdToLanguages < ActiveRecord::Migration[8.0]
  def change
    add_column :languages, :openai_assistant_id, :string
  end
end
