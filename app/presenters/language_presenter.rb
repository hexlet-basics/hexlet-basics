# frozen_string_literal: true

module LanguagePresenter
  def source_code_curl
    "#{ExternalLinks.source_code_curl}/exercises-#{slug}"
  end
end
