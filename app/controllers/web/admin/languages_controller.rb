# frozen_string_literal: true

class Web::Admin::LanguagesController < Web::Admin::ApplicationController

  def index
    @languages = Language.all
  end

  def new
    @language = Language.new
  end

  def create
    @language = Language.new(language_params)

    if @language.save
      redirect_to admin_languages_path
    else
      render :new
    end
  end

  private

  def language_params
    params.require(:language).permit(:slug)
  end
end
