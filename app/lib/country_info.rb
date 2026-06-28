# typed: false
# frozen_string_literal: true

# NOTE: dead code — `CountryInfo` has no callers and references `ISO3166`
# (the `countries` gem), which is not in the Gemfile, so `current_name` would
# raise NameError if invoked. Kept as typed:false pending removal.
class CountryInfo
  def self.current_name
    codes = { ru: "7", en: "1" }
    country = ISO3166::Country[codes.fetch(I18n.locale)]
    country.translations[I18n.locale.to_s]
  end
end
