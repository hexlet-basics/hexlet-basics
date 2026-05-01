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
    category = Admin::LanguageCategoryForm.new
    category.locale = I18n.locale
    landing_pages = Language::LandingPage.web.merge(Language.ordered)

    render inertia: true, props: {
      categoryDto: Language::CategoryCreateResource.new(category),
      qnaItems: Language::CategoryQnaItemResource.new([]),
      landingPagesForCategories: Language::LandingPageForListsResource.new(landing_pages)
    }
  end

  def edit
    category = Admin::LanguageCategoryForm.find(params[:id])
    category.locale = I18n.locale
    landing_pages = Language::LandingPage.web.merge(Language.ordered)

    render inertia: true, props: {
      categoryDto: Language::CategoryUpdateResource.new(category),
      qnaItems: Language::CategoryQnaItemResource.new(category.qna_items.order(:id)),
      landingPagesForCategories: Language::LandingPageForListsResource.new(landing_pages)
    }
  end

  def create
    category = Admin::LanguageCategoryForm.new(params[:data])
    category.locale = I18n.locale

    if category.save
      f(:success)
      redirect_to edit_admin_language_category_path(category), inertia: { errors: category.errors }
    else
      f(:error)
      redirect_to new_admin_language_category_url, inertia: { errors: category.errors }
    end
  end

  def update
    category = Admin::LanguageCategoryForm.find(params[:id])
    category.locale = I18n.locale

    if category.update(params[:data])
      f(:success)
    else
      f(:error)
    end

      redirect_to edit_admin_language_category_path(category), inertia: { errors: category.errors }
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
