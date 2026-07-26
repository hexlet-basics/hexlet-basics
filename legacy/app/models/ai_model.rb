# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: ai_models
#
#  id                :bigint           not null, primary key
#  capabilities      :jsonb
#  context_window    :integer
#  family            :string
#  knowledge_cutoff  :date
#  max_output_tokens :integer
#  metadata          :jsonb
#  modalities        :jsonb
#  model_created_at  :datetime
#  name              :string           not null
#  pricing           :jsonb
#  provider          :string           not null
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  model_id          :string           not null
#
# Indexes
#
#  index_ai_models_on_capabilities           (capabilities) USING gin
#  index_ai_models_on_family                 (family)
#  index_ai_models_on_modalities             (modalities) USING gin
#  index_ai_models_on_provider_and_model_id  (provider,model_id) UNIQUE
#
class AiModel < ApplicationRecord
  acts_as_model chats: :ai_chats
end
