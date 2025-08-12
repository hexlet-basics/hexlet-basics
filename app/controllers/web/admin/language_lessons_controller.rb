class Web::Admin::LanguageLessonsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Language::Lesson::Version::Info.current.with_locale.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      lessons: Language::LessonResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
