# frozen_string_literal: true

xml.instruct! :xml, version: '1.0', encoding: 'utf-8', standalone: :yes
xml.yml_catalog date: l(Time.current, format: :date_time) do
  xml.shop do
    xml.name t('.title')
    xml.company t('.company')
    xml.url root_url
    xml.email 'support@hexlet.io'
    xml.currencies do
      xml.currency id: 'RUR', rate: 1
    end
    xml.categories do
      xml.category t('.category'), id: 1
    end
    xml.offers do
      @languages.each do |language|
        xml.offer id: language.slug, available: true do
          xml.name language.current_version.name.downcase
          info = @infos_by_language[language.id]
          xml.description info.description.truncate(250, separator: ' ')
          xml.url language_url(language.slug, locale: AppHost.locale_for_url(info.locale))
          xml.categoryId 1
          xml.currencyId 'RUR'
          xml.price 0
          xml.available true
          xml.picture image_url("#{language.slug}.png")
        end
      end
    end
  end
end
