# frozen_string_literal: true

class Api::Languages::ApplicationController < Api::ApplicationController
  def resource_language
    @resource_language ||= Language.find(params[:language_id])
  end

  def resource_language_version
    @resource_language_version ||= resource_language.current_version
  end
end
