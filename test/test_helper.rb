ENV["RAILS_ENV"] ||= "test"
require_relative "../config/environment"
require "rails/test_help"

module ActiveSupport
  class TestCase
    # Run tests in parallel with specified workers
    parallelize(workers: :number_of_processors)

    # Setup all fixtures in test/fixtures/*.yml for all tests in alphabetical order.
    fixtures :all

    # Add more helper methods to be used by all tests here...
  end
end

# # frozen_string_literal: true
#
# ENV['RAILS_ENV'] ||= 'test'
# require_relative '../config/environment'
# require 'rails/test_help'
#
# # rubocop:disable Rails/I18nLocaleAssignment
# I18n.locale = ENV.fetch('RAILS_LOCALE', 'en').to_sym
# # rubocop:enable Rails/I18nLocaleAssignment
# Rails.application.routes.default_url_options[:locale] = AppHost.locale_for_url(I18n.locale)
#
# OmniAuth.config.test_mode = true
# OmniAuth.config.request_validation_phase = OmniAuth::AuthenticityTokenProtection.new(allow_if: ->(_env) { true })
#
# class ActiveSupport::TestCase
#   include FactoryBot::Syntax::Methods
#
#   parallelize(workers: :number_of_processors)
#
#   fixtures :all
#
#   def setup
#     queue_adapter.perform_enqueued_jobs = true
#     queue_adapter.perform_enqueued_at_jobs = true
#   end
# end
#
# class ActionDispatch::IntegrationTest
#   include AuthConcern
#
#   def sign_in_as(name)
#     user = users(name)
#
#     post session_url, params: { sign_in_form: { email: user.email, password: 'password' } }
#     user
#   end
# end
