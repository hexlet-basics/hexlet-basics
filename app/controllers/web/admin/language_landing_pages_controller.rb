class Web::Admin::LanguageLandingPagesController < Web::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Language::LandingPage.with_locale.joins(:language).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      landingPages: Language::LandingPageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    landing_page = Admin::LanguageLandingPageForm.new
    landing_page.locale = I18n.locale

    render inertia: true, props: {
      landingPageDto: Language::LandingPageCrudResource.new(landing_page)
    }
  end

  def edit
    landing_page = Admin::LanguageLandingPageForm.with_locale.find(params[:id])
    # versions = language.versions.limit(5).order(created_at: :desc)

    render inertia: true, props: {
      landingPageDto: Language::LandingPageCrudResource.new(landing_page)
      # courseVersions: Language::VersionResource.new(versions)
    }
  end

  def create
    landing_page = Admin::LanguageLandingPageForm.new(params[:language_landing_page])

    if landing_page.save
      f(:success)
      redirect_to admin_language_landing_pages_path
    else
      f(:error)
      redirect_to_inertia new_admin_language_landing_page_path, landing_page
    end
  end

  def update
    landing_page = Admin::LanguageLandingPageForm.with_locale.find(params[:id])
    # raise params.inspect

    if landing_page.update(params[:language_landing_page])
      f(:success)
    else
      # raise landing_page.qna_items.inspect
      raise landing_page.errors.full_messages.inspect
      f(:error)
    end
    redirect_to_inertia edit_admin_language_landing_page_path(landing_page), landing_page
  end
end
