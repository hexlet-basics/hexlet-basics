# typed: strict

class Web::Languages::ApplicationController < Web::ApplicationController
  sig { returns(Language) }
  def resource_language
    @resource_language ||= T.let(Language.find_by!(slug: params.expect(:language_id)), T.nilable(Language))
  end

  sig { returns(Language::LandingPage) }
  def resource_language_landing_page
    @resource_language_landing_page ||= T.let(resource_language.landing_pages.published.with_locale.find_by!(main: true), T.nilable(Language::LandingPage))
    # TODO: implement redirect logic
    # @resource_language_landing_page ||= resource_language.landing_pages.find_by!(locale: I18n.locale)
  end

  # alias_method :resource_language_landing_page, :landing_page
end
