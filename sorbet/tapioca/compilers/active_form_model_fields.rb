# typed: strict
# frozen_string_literal: true

require "tapioca/dsl"

module Tapioca
  module Dsl
    module Compilers
      # Generates RBI for ActiveFormModel::Virtual forms. `fields :a, :b` defines
      # `attr_accessor`s for each permitted arg; this compiler emits reader/writer
      # sigs for `_permitted_args` so typed forms can reference their fields.
      class ActiveFormModelFields < Compiler
        extend T::Sig

        ConstantType = type_member { { fixed: T::Class[T.anything] } }

        sig { override.returns(T::Enumerable[T::Module[T.anything]]) }
        def self.gather_constants
          all_classes.select do |klass|
            klass.include?(::ActiveFormModel::Virtual) && klass.respond_to?(:_permitted_args)
          end
        end

        sig { override.void }
        def decorate
          fields = T.unsafe(constant)._permitted_args
          return if fields.blank?

          root.create_path(constant) do |klass|
            fields.each do |field|
              klass.create_method(field.to_s, return_type: "T.untyped")
              klass.create_method(
                "#{field}=",
                parameters: [ create_param("value", type: "T.untyped") ],
                return_type: "T.untyped"
              )
            end
          end
        end
      end
    end
  end
end
