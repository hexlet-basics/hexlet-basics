# frozen_string_literal: true

module BrowserConcern
  extend ActiveSupport::Concern

  class_methods do
    def allow_modern_browsers
      before_action do
        unless modern_browser?
          # https://github.com/rails/rails/blob/4d42d34adda4011e770c18136966cb5df08c9220/actionpack/lib/action_controller/metal/allow_browser.rb#L57
          render file: Rails.root.join("public/406-unsupported-browser.html"), layout: false, status: :not_acceptable
        end
      end
    end
  end

  private

  # https://vite.dev/guide/build#browser-compatibility
  def modern_browser?
    [
      browser.chrome?(">= 87"),
      browser.safari?(">= 14"),
      browser.firefox?(">= 78"),
      browser.edge?(">= 88"),
      browser.opera?(">= 73"),
      browser.facebook? &&
        browser.safari_webapp_mode? &&
        browser.webkit_full_version.to_i >= 605
    ].any?
  end

  def mobile_browser?
    browser.device.mobile?
  end
end
