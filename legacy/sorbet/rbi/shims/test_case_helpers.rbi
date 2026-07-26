# typed: true

# Rails mixes these modules into the test base classes at boot (route helpers via
# the app's routes, ActiveJob/ActionMailer test helpers via railties, Capybara DSL
# into SystemTestCase). Tapioca doesn't run those initializers, so make the
# includes explicit here so tests can call `root_path`, `assert_enqueued_jobs`,
# `visit`, etc.
class ActionDispatch::IntegrationTest
  include GeneratedPathHelpersModule
  include GeneratedUrlHelpersModule
  include ActiveJob::TestHelper
  include ActionMailer::TestHelper
end

class ActionDispatch::SystemTestCase
  include GeneratedPathHelpersModule
  include GeneratedUrlHelpersModule
  include Capybara::DSL
  include Capybara::Minitest::Assertions
end
