class Api::FeedsController < Api::ApplicationController
  def yandex_courses
    I18n.with_locale :ru do
      landingPages = Language::LandingPage.web
        .includes([ { language: { cover_attachment: :blob } } ])
      categories = Language::Category.with_locale
      respond_to do |format|
        format.xml { render xml: YandexSearchFeedBuilder.build(landingPages, categories) }
      end
    end
  end
end
