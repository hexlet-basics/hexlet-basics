# frozen_string_literal: true

class Api::LanguagesController < Api::ApplicationController
  def index
    @languages = Language.with_progress(:completed).includes(:current_version)
    respond_with @languages
  end

  def show
    @language = Language.with_progress(:completed).find(params[:id])
    respond_with @language
  end
end
