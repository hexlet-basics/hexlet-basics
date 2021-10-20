# frozen_string_literal: true

class Api::LanguagesController < Api::ApplicationController
  def index
    @languages = Language.with_progress(:completed).includes(:current_version)
  end

  def show
    @language = Language.with_progress(:completed).find(params[:id])
  end
end
