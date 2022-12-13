# frozen_string_literal: true

# Be sure to restart your server when you modify this file.

Rails.application.config.session_store :cookie_store, key: '_hexlet_basics_session',
                                                      expire_after: 1.month,
                                                      domain: Rails.env.test? ? ENV.fetch('APP_HOST') : :all, # NOTE: :all not working on test env
                                                      same_site: :lax
