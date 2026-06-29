# typed: strict
# frozen_string_literal: true

require "tapioca/dsl"

# Rails mixes `ActiveRecord::TestFixtures` into `ActiveSupport::TestCase` via the
# `active_record.test_fixtures` railtie initializer. Tapioca doesn't run
# initializers during gem generation, so the include is absent from the gem RBI
# and Sorbet can't find the `fixtures` class method (it lives in
# `ActiveRecord::TestFixtures::ClassMethods`, added via `mixes_in_class_methods`).
#
# This DSL compiler makes the include explicit in the RBI for
# `ActiveSupport::TestCase` so `fixtures` (and the generated accessors) resolve.
module Tapioca
  module Dsl
    module Compilers
      class ActiveSupportTestCaseTestFixtures < Compiler
        extend T::Sig

        ConstantType = type_member { { fixed: T.class_of(::ActiveSupport::TestCase) } }

        sig { override.returns(T::Enumerable[T::Class[T.anything]]) }
        def self.gather_constants
          return [] unless defined?(::Rails.application) && ::Rails.application
          return [] unless defined?(::ActiveRecord::TestFixtures)

          [ ::ActiveSupport::TestCase ]
        end

        sig { override.void }
        def decorate
          root.create_path(constant) do |scope|
            scope.create_include("::ActiveRecord::TestFixtures")
          end
        end
      end
    end
  end
end
