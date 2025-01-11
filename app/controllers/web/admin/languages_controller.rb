# frozen_string_literal: true

class Web::Admin::LanguagesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("s" => "created_at desc")
    search = Language::Version::Info.with_locale.includes([ language: :current_version ]).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      languages: LanguageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    @language = Language.new

    render inertia: true, props: {
    }
  end

  def edit
    @language = Language.find(params[:id])

    render inertia: true, props: {
    }
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
    params.require(:language).permit(:slug, :progress, :learn_as, :order, :category_id)
  end
end
