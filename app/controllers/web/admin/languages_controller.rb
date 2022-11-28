# frozen_string_literal: true

class Web::Admin::LanguagesController < Web::Admin::ApplicationController
  def index
    q = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')
    @search = Language.ransack(q)
    @languages = @search.result
  end

  def new
    @language = Language.new
  end

  def edit
    @language = Language.find(params[:id])
  end

  def create
    @language = Language.new(language_params)

    if @language.save
      f(:success)
      redirect_to admin_languages_path
    else
      f(:error)
      render :new
    end
  end

  def update
    @language = Language.find(params[:id])

    if @language.update(language_params)
      f(:success)
      redirect_to admin_languages_path
    else
      f(:error)
      render :edit
    end
  end

  private

  def language_params
    params.require(:language).permit(:slug, :progress, :learn_as, :order)
  end
end
