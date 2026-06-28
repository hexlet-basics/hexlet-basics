# typed: true

class Web::Admin::LanguageLandingPagesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "id", "so" => "desc", fields: { "state_eq" => "published" })
    search = Language::LandingPage.with_locale.joins(:language).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      landingPages: Language::LandingPageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    landing_page = Language::LandingPage.new
    landing_pages = Language::LandingPage.published

    languages = Language.all

    render inertia: true, props: {
      landingPageDto: Language::LandingPageCreateResource.new(landing_page),
      landingPages: Language::LandingPageResource.new(landing_pages),
      languages: LanguageResource.new(languages),
      qnaItems: Language::LandingPageQnaItemResource.new([])
    }
  end

  def edit
    landing_page = Language::LandingPage.with_locale.find(params[:id])
    languages = Language.all
    landing_pages = Language::LandingPage.published

    render inertia: true, props: {
      landingPageDto: Language::LandingPageUpdateResource.new(landing_page),
      landingPages: Language::LandingPageResource.new(landing_pages),
      languages: LanguageResource.new(languages),
      qnaItems: Language::LandingPageQnaItemResource.new(landing_page.qna_items.order(:id))
    }
  end

  def create
    struct = ApplicationParamsStruct.from_params(LandingPageStruct, params.require(:data))
    result = LandingPageService.create(struct, locale: I18n.locale.to_s, outcomes_image: params.dig(:data, :outcomes_image))

    case result
    when Typed::Success
      f(:success)
      redirect_to admin_language_landing_pages_path
    when Typed::Failure
      f(:error)
      redirect_to new_admin_language_landing_page_path, inertia: { errors: result.error.errors }
    end
  end

  def update
    struct = ApplicationParamsStruct.from_params(LandingPageStruct, params.require(:data))
    result = LandingPageService.update(params[:id], struct, outcomes_image: params.dig(:data, :outcomes_image))

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_language_landing_page_path(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_language_landing_page_path(result.error), inertia: { errors: result.error.errors }
    end
  end
end
