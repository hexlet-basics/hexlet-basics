# typed: strict

class Web::Admin::LanguageLessonMembersController < Web::Admin::ApplicationController
  STAFF_RESOURCE = StaffMember::Role::Permission::Resource::LanguageLessonMembers

  sig { void }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Language::Lesson::Member
      .includes([ [ lesson: :infos ], :language ])
      .ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      languageLessonMembers: Language::Lesson::MemberResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
