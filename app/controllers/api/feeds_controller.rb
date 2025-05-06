class Api::FeedsController < Api::ApplicationController
  def yandex_courses
    I18n.with_locale :ru do
      landingPages = Language::LandingPage.web.where(listed: true).where(main: true)
      categories = Language::Category.with_locale
      respond_to do |format|
        format.xml { render xml: YandexSearchFeedBuilder.build(landingPages, categories) }
      end
    end
  end
end
