# frozen_string_literal: true

class Web::Languages::ApplicationController < Web::ApplicationController
  helper_method :resource_language

  def resource_language
    @resource_language ||= Language.find_by!(slug: params[:language_id])
  end
end
