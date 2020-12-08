# frozen_string_literal: true

ENV['RAILS_ENV'] ||= 'test'
require_relative '../config/environment'
require 'rails/test_help'

OmniAuth.config.test_mode = true
class ActiveSupport::TestCase
  include FactoryBot::Syntax::Methods

  parallelize(workers: :number_of_processors)

  fixtures :all

  def setup
    queue_adapter.perform_enqueued_jobs = true
    queue_adapter.perform_enqueued_at_jobs = true
  end
end

class ActionDispatch::IntegrationTest
  include AuthManagment

  def sign_in_as(name)
    user = users(name)

    post session_path, params: { sign_in: { email: user.email, password: 'password' } }
    user
  end
end
