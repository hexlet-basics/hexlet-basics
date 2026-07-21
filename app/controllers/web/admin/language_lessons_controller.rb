# typed: strict

class Web::Admin::LanguageLessonsController < Web::Admin::ApplicationController
  STAFF_RESOURCE = StaffMember::Role::Permission::Resource::LanguageLessons

  sig { void }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Language::Lesson::Version::Info.current.with_locale.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      lessons: Language::LessonResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  sig { void }
  def review
    lesson = Language::Lesson.find(params[:id])
    lesson.infos.find_each do |info|
      ReviewLessonJob.perform_later(info.id)
    end

    f(:success)

    redirect_to admin_language_lessons_path
  end
end
