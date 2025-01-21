# frozen_string_literal: true

class Web::ErrorsController < Web::ApplicationController
  # NOTE: for 404 page locale from route params is undefined
  around_action :use_locale # , only: :not_found

  def show
    set_meta_tags title: t(".base")

    render inertia: true, props: {
      code: params[:code],
      message: Rack::Utils::HTTP_STATUS_CODES[params[:code].to_i]
    }
  end

  private

  def use_locale(&)
    locale = URI.parse(request.original_url).path.split("/").second || ""
    locale = I18n.default_locale unless I18n.available_locales.include?(locale.to_sym)

    I18n.with_locale(locale, &)
  end
end
