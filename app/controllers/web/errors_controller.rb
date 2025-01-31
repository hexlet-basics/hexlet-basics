# frozen_string_literal: true

class Web::ErrorsController < Web::ApplicationController
  # NOTE: for 404 page locale from route params is undefined
  around_action :use_locale # , only: :not_found

  def show
    code = params[:code].to_i
    header = t(".codes.#{code}.header", default: t(".codes.other.header"))
    description = t(".codes.#{code}.description", default: t(".codes.other.header"))

    seo_tags = {
      title: header,
      description: description
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      code:,
      header:,
      description:
    }
  end

  private

  def use_locale(&)
    locale = URI.parse(request.original_url).path.split("/").second || ""
    locale = I18n.default_locale unless I18n.available_locales.include?(locale.to_sym)

    I18n.with_locale(locale, &)
  end
end
