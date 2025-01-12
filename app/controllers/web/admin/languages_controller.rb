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
    render inertia: true, props: {
    }
  end

  def edit
    language_info = Language::Version::Info.find(params[:id])

    render inertia: true, props: {
      language: LanguageResource.new(language_info)
    }
  end

  def create
    language = Language.new(language_params)

    if language.save
      f(:success)
      redirect_to admin_languages_path
    else
      f(:error)
      redirect_to_inertia new_admin_language_path, language
    end
  end

  def update
    language_info = Language::Version::Info.find(params[:id])

    if language_info.update(language_params)
      f(:success)
      redirect_to admin_languages_path
    else
      f(:error)
      redirect_to_inertia edit_admin_language_path(language), language
    end
  end

  private

  def language_params
    params.require(:language).permit(:slug, :progress, :learn_as, :order, :category_id)
  end
end
