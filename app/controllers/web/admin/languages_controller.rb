# typed: strict

class Web::Admin::LanguagesController < Web::Admin::ApplicationController
  STAFF_RESOURCE = StaffMember::Role::Permission::Resource::Languages

  sig { returns(T.untyped) }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Language.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      courses: LanguageResource.new(records),
      # mainLandingPagesByLanguageId: main_landing_pages_by_landing_id
      grid: GridResource.new(grid_params(pagy))
    }
  end

  sig { returns(T.untyped) }
  def new
    language = Language.new
    render inertia: true, props: {
      courseDto: LanguageCreateResource.new(language)
    }
  end

  sig { returns(T.untyped) }
  def edit
    language = Language.find(params[:id])
    versions = language.versions.limit(5).order(created_at: :desc)
    landing_page = language.landing_pages.published.find_by(locale: I18n.locale, main: true)

    render inertia: true, props: {
      courseDto: LanguageUpdateResource.new(language),
      landingPage: landing_page && Language::LandingPageResource.new(landing_page),
      courseVersions: Language::VersionResource.new(versions)
    }
  end

  sig { returns(T.untyped) }
  def create
    struct = ApplicationParamsStruct.from_params(LanguageStruct, params.require(:data))
    result = LanguageService.create(struct, cover: params.dig(:data, :cover))

    case result
    when Typed::Success
      f(:success)
      redirect_to admin_languages_path
    when Typed::Failure
      f(:error)
      redirect_to new_admin_language_path, inertia: { errors: result.error.errors }
    end
  end

  sig { returns(T.untyped) }
  def review
    language = Language.find(params[:id])
    language.current_lesson_infos.find_each do |info|
      ReviewLessonJob.perform_later(info.id)
    end

    f(:success)

    redirect_to admin_languages_path
  end

  sig { returns(T.untyped) }
  def update
    struct = ApplicationParamsStruct.from_params(LanguageStruct, params.require(:data))
    result = LanguageService.update(params[:id], struct, cover: params.dig(:data, :cover))

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_language_path(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_language_path(result.error), inertia: { errors: result.error.errors }
    end
  end
end
