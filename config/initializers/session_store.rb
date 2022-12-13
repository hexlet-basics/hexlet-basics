# frozen_string_literal: true

# Be sure to restart your server when you modify this file.

Rails.application.config.session_store :cookie_store, key: '_hexlet_basics_session',
                                                      expire_after: 1.month,
                                                      domain: :all,
                                                      same_site: :lax
