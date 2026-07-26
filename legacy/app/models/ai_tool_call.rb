# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: ai_tool_calls
#
#  id                :bigint           not null, primary key
#  arguments         :jsonb
#  name              :string           not null
#  thought_signature :text
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  ai_message_id     :bigint           not null
#  tool_call_id      :string           not null
#
# Indexes
#
#  index_ai_tool_calls_on_ai_message_id  (ai_message_id)
#  index_ai_tool_calls_on_name           (name)
#  index_ai_tool_calls_on_tool_call_id   (tool_call_id) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (ai_message_id => ai_messages.id)
#
class AiToolCall < ApplicationRecord
  acts_as_tool_call message: :ai_message
end
