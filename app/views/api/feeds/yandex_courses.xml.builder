# https://cachev2-kiv-06.cdn.yandex.net/download.cdn.yandex.net/from/yandex.ru/support/ru/webmaster/files/education.xml?lid=307
xml.instruct! :xml, version: "1.0", encoding: "UTF-8"
xml.yml_catalog(date: Time.now.strftime("%Y-%m-%d %H:%M")) do
  xml.shop do
    xml.name "Название вашей школы"
    xml.company "Название компании"
    xml.url root_url
    xml.email "support@example.com"
    xml.picture "https://example.com/logo.png"
    xml.description "Описание вашей образовательной платформы"

    # xml.currencies do
    #   xml.currency(id: "RUR", rate: "0")
    # end

    xml.offers do
      @landingPages.ach do |lp|
        xml.offer(id: lp.language.id) do
          xml.name lp.header
          xml.url language_url(lp.slug)
          # xml.categoryId course.category_id
          xml.price 0
          xml.currencyId "RUR"
          xml.param(name: "Есть бесплатная часть") { xml.text true }
          xml.param(name: "Продолжительность", unit: "час") { xml.text lp.language.duration }
          xml.param(name: "Формат обучения") { xml.text "Самостоятельно" }
          xml.param(name: "Сложность") { xml.text "Для новичков" }
          xml.param(name: "Тип обучения") { xml.text "Курс" }
          xml.description lp.description
          xml.picture rails_representation_url(lp.language.cover.variant(:list)) if lp.language.cover.attached?
        end
      end
    end
  end
end
