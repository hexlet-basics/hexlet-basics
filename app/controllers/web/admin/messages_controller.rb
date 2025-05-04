class Web::Admin::MessagesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    # raise q.inspect
    search = Language::Lesson::Member::Message
      .joins(:language, :language_lesson)
      .includes([ :language_lesson_member ]).ransack(q)
    # raise search.result.to_sql
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      messages: Language::Lesson::Member::MessageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def by_lesson_member
    scope = Language::Lesson::Member.order(id: :desc)
    pagy, records = pagy(scope)

    render inertia: true, props: {
      members: Language::Lesson::MemberResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
