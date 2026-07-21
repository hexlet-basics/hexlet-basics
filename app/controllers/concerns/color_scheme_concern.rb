# typed: strict
# frozen_string_literal: true

module ColorSchemeConcern
  COLOR_SCHEME_COOKIE_KEY = "codebasics-color-scheme"
  VALID_COLOR_SCHEMES = %w[light dark].freeze

  extend ActiveSupport::Concern

  extend T::Sig

  extend T::Helpers

  requires_ancestor { ApplicationController }

  included do
    T.bind(self, T.class_of(ApplicationController))
    helper_method :current_color_scheme
  end

  sig { returns(String) }
  def current_color_scheme
    color_scheme = cookies[COLOR_SCHEME_COOKIE_KEY]

    VALID_COLOR_SCHEMES.include?(color_scheme) ? color_scheme : "light"
  end
end
