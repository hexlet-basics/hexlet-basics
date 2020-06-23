# frozen_string_literal: true

class Web::Admin::Languages::ApplicationController < Web::Admin::ApplicationController
  helper_method :resource_language

  def resource_language
    @resource_language ||= Language.find(params[:language_id])
  end
end
