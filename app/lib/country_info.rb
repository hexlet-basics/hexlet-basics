# frozen_string_literal: true

class CountryInfo
  def self.current_name
    codes = { ru: '7', en: '1' }
    country = ISO3166::Country[codes.fetch(I18n.locale)]
    country.translations[I18n.locale.to_s]
  end
end
