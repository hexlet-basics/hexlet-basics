# typed: strict

class Web::Admin::LanguageLessonReviewsController < Web::Admin::ApplicationController
  STAFF_RESOURCE = StaffMember::Role::Permission::Resource::LanguageLessonReviews

  sig { void }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Language::Lesson::Review.where(locale: I18n.locale).with_summary.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      reviews: Language::Lesson::ReviewResource.new(records),
      languages: LanguageResource.new(Language.web.order(slug: :asc)),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
