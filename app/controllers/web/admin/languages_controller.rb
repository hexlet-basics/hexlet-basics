# frozen_string_literal: true

class Web::Admin::LanguagesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Language.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      courses: LanguageResource.new(records),
      # mainLandingPagesByLanguageId: main_landing_pages_by_landing_id
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    language = Admin::LanguageForm.new
    render inertia: true, props: {
      courseDto: LanguageCrudResource.new(language)
    }
  end

  def edit
    language = Admin::LanguageForm.find(params[:id])
    versions = language.versions.limit(5).order(created_at: :desc)
    landing_page = language.landing_pages.published.find_by!(locale: I18n.locale, main: true)

    render inertia: true, props: {
      courseDto: LanguageCrudResource.new(language),
      landingPage: Language::LandingPageResource.new(landing_page),
      courseVersions: Language::VersionResource.new(versions)
    }
  end

  def create
    language = Admin::LanguageForm.new(params[:language])

    if language.save
      f(:success)
      redirect_to admin_languages_path
    else
      f(:error)
      redirect_to_inertia new_admin_language_path, language
    end
  end

  def update
    language = Admin::LanguageForm.find(params[:id])

    if language.update(params[:language])
      f(:success)
    else
      f(:error)
    end
    redirect_to_inertia edit_admin_language_path(language), language
  end
end
