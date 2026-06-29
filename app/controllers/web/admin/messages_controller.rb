# typed: strict

class Web::Admin::MessagesController < Web::Admin::ApplicationController
  sig { returns(T.untyped) }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = AiMessage
      .includes(ai_chat: { language_lesson_member: { lesson: %i[language versions] } })
      .ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      messages: AiMessageResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  sig { returns(T.untyped) }
  def by_lesson_member
    scope = Language::Lesson::Member.order(id: :desc)
    pagy, records = T.unsafe(self).pagy(scope)

    render inertia: true, props: {
      members: Language::Lesson::MemberResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
