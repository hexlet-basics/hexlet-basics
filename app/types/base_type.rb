# frozen_string_literal: true

module BaseType
  extend ActiveSupport::Concern

  module ClassMethods
    delegate :model_name, to: :superclass
  end
end
