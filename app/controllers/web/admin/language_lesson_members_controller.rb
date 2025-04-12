class Web::Admin::LanguageLessonMembersController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Language::Lesson::Member
      .order(id: :desc).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      languageLessonMembers: Language::Lesson::MemberResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
