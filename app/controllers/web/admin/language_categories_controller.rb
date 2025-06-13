class Web::Admin::LanguageCategoriesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "desc")
    search = Language::Category.with_locale.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      categories: Language::CategoryResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    category = Admin::LanguageCategoryForm.new
    category.locale = I18n.locale

    render inertia: true, props: {
      categoryDto: Language::CategoryCrudResource.new(category)
    }
  end

  def edit
    category = Admin::LanguageCategoryForm.find(params[:id])
    category.locale = I18n.locale

    render inertia: true, props: {
      categoryDto: Language::CategoryCrudResource.new(category)
    }
  end

  def create
    category = Admin::LanguageCategoryForm.new(params[:language_category])
    category.locale = I18n.locale

    if category.save
      f(:success)
      redirect_to_inertia edit_admin_language_category_path(category), category
    else
      f(:error)
      redirect_to_inertia new_admin_language_category_url, category
    end
  end

  def update
    category = Admin::LanguageCategoryForm.find(params[:id])
    category.locale = I18n.locale

    if category.update(params[:language_category])
      f(:success)
    else
      f(:error)
    end

      redirect_to_inertia edit_admin_language_category_path(category), category
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
