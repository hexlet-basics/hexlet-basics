# frozen_string_literal: true

class Web::ErrorsController < Web::ApplicationController
  before_action do
    set_meta_tags title: t(".base")
  end

  # NOTE: for 404 page locale from route params is undefined
  around_action :use_locale, only: :not_found

  def forbidden
    respond_to do |format|
      format.html { render status: :forbidden }
      format.all  { render nothing: true, status: :forbidden }
    end
  end

  def not_found
    respond_to do |format|
      format.html { render status: :not_found }
      format.all  { head :not_found }
    end
  end

  def server_error
    respond_to do |format|
      format.html { render status: :internal_server_error }
      format.all  { render nothing: true, status: :internal_server_error }
    end
  end

  private

  def use_locale(&)
    locale = URI.parse(request.original_url).path.split("/").second || ""
    locale = I18n.default_locale unless I18n.available_locales.include?(locale.to_sym)

    I18n.with_locale(locale, &)
  end
end
