# typed: strict
# frozen_string_literal: true

require "tapioca/dsl"

module Tapioca
  module Dsl
    module Compilers
      # Adds `acts_as_taggable_on` class method to ApplicationRecord, plus
      # `tagged_with` / `tag_counts_on` to ActiveRecord::Relation so Sorbet
      # sees them on the per-model `*::PrivateRelation` / `PrivateAssociationRelation`
      # / `PrivateCollectionProxy` classes (all of which extend Relation).
      class ActsAsTaggableOn < Compiler
        extend T::Sig

        ConstantType = type_member { { fixed: T.class_of(::ApplicationRecord) } }

        sig { override.returns(T::Enumerable[T::Class[T.anything]]) }
        def self.gather_constants
          [ ::ApplicationRecord ]
        end

        sig { override.void }
        def decorate
          root.create_path(constant) do |scope|
            scope.create_method(
              "acts_as_taggable_on",
              class_method: true,
              parameters: [
                create_rest_param("contexts", type: "T.untyped")
              ],
              return_type: "void"
            )
          end

          root.create_path(::ActiveRecord::Relation) do |scope|
            scope.create_method(
              "tagged_with",
              parameters: [
                create_param("tags", type: "T.untyped"),
                create_opt_param("options", type: "T.untyped", default: "{}")
              ],
              return_type: "ActiveRecord::Relation"
            )

            scope.create_method(
              "tag_counts_on",
              parameters: [
                create_param("context", type: "T.untyped"),
                create_opt_param("options", type: "T.untyped", default: "{}")
              ],
              return_type: "ActiveRecord::Relation"
            )
          end
        end
      end
    end
  end
end
