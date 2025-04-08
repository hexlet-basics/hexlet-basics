module YandexSearchFeedBuilder
  def self.build(landingPages, categories)
    urls = Rails.application.routes.url_helpers
    builder = Nokogiri::XML::Builder.new do |xml|
      # xml.instruct! :xml, version: "1.0", encoding: "UTF-8"
      xml.yml_catalog(date: Time.now.strftime("%Y-%m-%d %H:%M")) do
        xml.shop do
          xml.name I18n.t("common.organization.name")
          xml.company I18n.t("common.organization.legal_name")
          xml.url urls.root_url(suffix: :ru)
          xml.email I18n.t("common.organization.email")
          xml.description I18n.t("common.organization.description")
          xml.picture "https://code-basics.com/images/logo.png"

          # xml.currencies do
          #   xml.currency(id: "RUR", rate: "0")
          # end

          xml.sets do
            categories.each do |category|
              xml.set(id: category.id) do
                xml.name category.name
                xml.url urls.language_category_url(category.slug, suffix: :ru)
              end
            end
          end

          xml.offers do
            landingPages.each do |lp|
              module_infos = lp.language.current_module_infos.with_locale

              if module_infos.size < 3
                next
              end

              utm_params = {
                utm_source: :yandex,
                utm_medium: :organic,
                utm_content: :feed_search,
                utm_campaign: "cb-#{lp.slug}",
                utm_term: "page_#{lp.slug}"
              }

              xml.offer(id: lp.language.id) do
                xml.name lp.header
                xml.url urls.language_url(lp.slug, suffix: :ru, **utm_params)
                # TODO добавить на сайт в базу и сделать прямую привязку курсов к категориям
                # https://yastatic.net/s3/doc-binary/src/support/products/education_rubricator.xml
                xml.categoryId 101
                xml.price 0
                xml.currencyId "RUR"
                xml.send "set-ids", lp.language_category_id
                xml.param(name: "Продолжительность", unit: "месяц") { xml.text 1 }
                xml.param(name: "Формат обучения") { xml.text "Самостоятельно" }
                xml.param(name: "Сложность") { xml.text "Для новичков" }
                xml.param(name: "Тип обучения") { xml.text "Курс" }
                xml.param(name: "Результат обучения") { xml.text "Сертификат" }
                xml.param(name: "Есть бесплатная часть") { xml.text true }
                xml.param(name: "Есть текстовые уроки") { xml.text true }
                xml.param(name: "Есть тренажеры") { xml.text true }
                xml.param(name: "Есть сообщество") { xml.text true }
                xml.description lp.description
                if lp.language.cover.attached?
                  xml.picture urls.rails_representation_url(lp.language.cover.variant(:list))
                end
                module_infos.each_with_index do |mi, i|
                  xml.param(name: "План", order: "#{i + 1}", unit: mi.name, hours: 5) do
                    xml.cdata <<-CDATA
                    #{mi.description}
                    CDATA
                  end
                end
              end
            end
          end
        end
      end
    end

    builder
  end
end
