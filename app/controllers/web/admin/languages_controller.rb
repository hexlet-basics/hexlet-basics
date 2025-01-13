# frozen_string_literal: true

class Web::Admin::LanguagesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("s" => "created_at desc")
    search = Language.includes([ :current_version ]).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      originalCourses: OriginalLanguageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    language = Language.new
    render inertia: true, props: {
      originalCourse: OriginalLanguageResource.new(language)
    }
  end

  def edit
    language = Language.find(params[:id])

    render inertia: true, props: {
      originalCourse: OriginalLanguageResource.new(language)
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
    language = Language.find(params[:id])

    if language.update(language_params)
      f(:success)
    else
      f(:error)
    end
    redirect_to_inertia edit_admin_language_path(language), language
  end

  private

  def language_params
    params.require(:language).permit(:slug, :progress, :learn_as, :order, :category_id)
  end
end
