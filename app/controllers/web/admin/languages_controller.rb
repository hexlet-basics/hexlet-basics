# frozen_string_literal: true

class Web::Admin::LanguagesController < Web::Admin::ApplicationController
  def index
    @languages = Language.all
  end

  def new
    @language = Language.new
  end

  def edit
    @language = Language.find(params[:id])
  end

  def update
    @language = Language.find(params[:id])

    if @language.update(language_params)
      redirect_to admin_languages_path
    else
      render :edit
    end
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
    params.require(:language).permit(:slug, :state_event)
  end
end
