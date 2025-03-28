class Web::Admin::MessagesController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Language::Lesson::Member::Message.user_role.includes([ :language_lesson_member ]).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      messages: Language::Lesson::Member::MessageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
