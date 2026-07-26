# typed: strict
# frozen_string_literal: true

require "tapioca/dsl"

module Tapioca
  module Dsl
    module Compilers
      # ransack mixes `ransack` into ActiveRecord relations at runtime; tapioca's
      # generated *::PrivateRelation classes don't get it. Add it for app models.
      class RansackRelations < Compiler
        extend T::Sig

        ConstantType = type_member { { fixed: T.class_of(::ApplicationRecord) } }

        RELATION_CLASS_NAMES = %w[PrivateRelation PrivateAssociationRelation PrivateCollectionProxy].freeze

        sig { override.returns(T::Enumerable[T::Class[T.anything]]) }
        def self.gather_constants
          all_classes.select do |klass|
            next false unless klass < ::ApplicationRecord
            next false if klass.abstract_class?

            source = Object.const_source_location(T.must(klass.name))&.first
            source&.include?("/app/models/") || false
          rescue ArgumentError, TypeError
            false
          end
        end

        sig { override.void }
        def decorate
          return if constant.name.nil?

          root.create_path(constant) do |model_scope|
            RELATION_CLASS_NAMES.each do |suffix|
              relation_class = model_scope.create_class(suffix)
              relation_class.create_method(
                "ransack",
                parameters: [
                  create_opt_param("params", type: "T.untyped", default: "nil"),
                  create_opt_param("options", type: "T.untyped", default: "nil")
                ],
                return_type: "T.untyped"
              )
            end
          end
        end
      end
    end
  end
end
