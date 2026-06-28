# typed: true

class Web::Admin::LanguageCategoriesController < Web::Admin::ApplicationController
  def index
    default_params = { "sf" => "id", "so" => "desc" }
    q = ransack_params(default_params)
    search = Language::Category.with_locale.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      categories: Language::CategoryResource.new(records),
      grid: GridResource.new(grid_params(pagy, default_params))
    }
  end

  def new
    category = Language::Category.new
    landing_pages = Language::LandingPage.web.merge(Language.ordered)

    render inertia: true, props: {
      categoryDto: Language::CategoryCreateResource.new(category),
      qnaItems: Language::CategoryQnaItemResource.new([]),
      landingPagesForCategories: Language::LandingPageForListsResource.new(landing_pages)
    }
  end

  def edit
    category = Language::Category.find(params[:id])
    landing_pages = Language::LandingPage.web.merge(Language.ordered)

    render inertia: true, props: {
      categoryDto: Language::CategoryUpdateResource.new(category),
      qnaItems: Language::CategoryQnaItemResource.new(category.qna_items.order(:id)),
      landingPagesForCategories: Language::LandingPageForListsResource.new(landing_pages)
    }
  end

  def create
    struct = ApplicationParamsStruct.from_params(LanguageCategoryStruct, params.require(:data))
    result = LanguageCategoryService.create(struct, locale: I18n.locale.to_s)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_language_category_path(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to new_admin_language_category_url, inertia: { errors: result.error.errors }
    end
  end

  def update
    struct = ApplicationParamsStruct.from_params(LanguageCategoryStruct, params.require(:data))
    result = LanguageCategoryService.update(params[:id], struct, locale: I18n.locale.to_s)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_language_category_path(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_language_category_path(result.error), inertia: { errors: result.error.errors }
    end
  end

  def destroy
    category = Language::Category.find(params[:id])

    # if category.language_landing_pages.any?
    #   f(:error)
    #   redirect_to admin_language_categories_path
    #   return
    # end

    if category.destroy
      f(:success)
    else
      f(:error)
    end

    redirect_to admin_language_categories_path
  end
end
