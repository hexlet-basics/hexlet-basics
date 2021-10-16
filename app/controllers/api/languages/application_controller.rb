# frozen_string_literal: true

class Api::Languages::ApplicationController < Api::ApplicationController
  helper_method :resource_language

  def resource_language
    @resource_langauge = Language.find(params[:language_id])
  end

  def resource_language_version
    @resource_langauge_version = resource_language.current_version
  end
end
