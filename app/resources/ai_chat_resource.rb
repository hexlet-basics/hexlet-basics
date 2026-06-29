# typed: strict

class AiChatResource < ApplicationResource
  typelize_from AiChat

  attributes :id, :created_at
end
