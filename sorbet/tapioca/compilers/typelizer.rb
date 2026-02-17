# typed: strict
# frozen_string_literal: true

require "tapioca/dsl"

module Tapioca
  module Dsl
    module Compilers
      # Generates RBI for ApplicationResource subclasses with typelize_from
      # to properly type the `object` method and `attribute` block parameter.
      class Typelizer < Compiler
        extend T::Sig

        ConstantType = type_member { { fixed: T.class_of(::ApplicationResource) } }

        sig { override.returns(T::Enumerable[Module]) }
        def self.gather_constants
          all_classes.select do |klass|
            klass < ::ApplicationResource && klass.respond_to?(:_typelizer_model_name)
          end
        end

        sig { override.void }
        def decorate
          return unless constant.respond_to?(:_typelizer_model_name)

          model_class = constant._typelizer_model_name
          return unless model_class

          model_name = model_class.name

          root.create_path(constant) do |klass|
            # Type the object method
            klass.create_method(
              "object",
              return_type: model_name
            )

            # Type the attribute method with typed block
            block_type = "T.proc.bind(#{constant.name}).params(it: #{model_name}).returns(T.untyped)"
            klass.create_method(
              "attribute",
              class_method: true,
              parameters: [
                create_opt_param("name", type: "T.any(Symbol, String)", default: "nil"),
                create_kw_rest_param("options", type: "T.untyped"),
                create_block_param("block", type: block_type)
              ],
              return_type: "void"
            )
          end
        end
      end
    end
  end
end
