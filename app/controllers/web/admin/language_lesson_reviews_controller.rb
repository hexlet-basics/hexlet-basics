class Web::Admin::LanguageLessonReviewsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Language::Lesson::Review.where(locale: I18n.locale).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      reviews: Language::Lesson::ReviewResource.new(records),
      languages: LanguageResource.new(Language.ordered),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
