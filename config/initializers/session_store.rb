# frozen_string_literal: true

Rails.application.config.session_store :cookie_store, key: COOKIE_STORE_KEY,
                                                      expire_after: 1.month
