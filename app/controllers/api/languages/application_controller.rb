# frozen_string_literal: true

class Api::Languages::ApplicationController < Api::ApplicationController
  helper_method :resource_language

  def resource_language
    @resource_language ||= Language.find(params[:language_id])
  end
end
