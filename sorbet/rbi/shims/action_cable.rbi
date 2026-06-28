# typed: true

# ActionCable::Channel::Base, Connection::Base and Channel::TestCase are
# autoloaded and not present when `tapioca gem` reflects actioncable, so they
# land in todo.rbi stubbed as modules — which trips Sorbet 5067 ("super class
# does not derive from Class") for our channels/tests that subclass them.
# Shim them as classes so the inheritance chain resolves; shims survive
# tapioca regeneration.
module ActionCable
  module Channel
    class Base; end
    class TestCase < ::ActiveSupport::TestCase; end
  end

  module Connection
    class Base; end
  end
end
