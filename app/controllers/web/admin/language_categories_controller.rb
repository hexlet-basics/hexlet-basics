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
    category = Admin::Language::CategoryForm.new
    category.creator = current_user

    render inertia: true, props: {
      cartegoryDto: Language::CategoryCrudResource.new(category)
    }
  end

  def edit
    category = Admin::Language::CategoryForm.find(params[:id])

    render inertia: true, props: {
      cartegoryDto: Language::CategoryCrudResource.new(category)
    }
  end

  def create
    category = Admin::Language::CategoryForm.new(params[:category])
    category.locale = I18n.locale
    category.creator = current_user

    if category.save
      f(:success)
      redirect_to_inertia edit_admin_category_path(category), category
    else
      f(:error)
      redirect_to_inertia new_admin_category_url, category
    end
  end

  def update
    category = Admin::Language::CategoryForm.find(params[:id])
    category.locale = I18n.locale

    if category.update(params[:category])
      f(:success)
    else
      f(:error)
    end

      redirect_to_inertia edit_admin_category_path(category), category
  end
end
