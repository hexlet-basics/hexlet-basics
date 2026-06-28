# typed: strict

class Web::Languages::ApplicationController < Web::ApplicationController
  sig { returns(T.untyped) }
  def resource_language
    @resource_language ||= T.let(Language.find_by!(slug: params[:language_id]), T.untyped)
  end

  sig { returns(T.untyped) }
  def resource_language_landing_page
    @resource_language_landing_page ||= T.let(resource_language.landing_pages.published.with_locale.find_by!(main: true), T.untyped)
    # TODO: implement redirect logic
    # @resource_language_landing_page ||= resource_language.landing_pages.find_by!(locale: I18n.locale)
  end

  # alias_method :resource_language_landing_page, :landing_page
end
