class Web::Admin::LanguageCategoriesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Language::Category.with_locale.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      categories: Language::CategoryResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    category = Admin::LanguageCategoryForm.new
    category.creator = current_user
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
    category.creator = current_user

    if category.save
      f(:success)
      redirect_to_inertia edit_admin_language_category_path(category), category
    else
      f(:error)
      redirect_to_inertia new_admin_category_url, category
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
end
