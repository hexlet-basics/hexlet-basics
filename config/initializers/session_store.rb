# frozen_string_literal: true

Rails.application.config.session_store :cookie_store, key: "_hexlet_basics_session_2",
                                                      expire_after: 1.month
