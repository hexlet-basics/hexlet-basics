class Api::FeedsController < Api::ApplicationController
  def yandex_courses
    @landingPages = Language::LandingPage.with_locale
    respond_to do |format|
      format.xml { render layout: false }
    end
  end
end
