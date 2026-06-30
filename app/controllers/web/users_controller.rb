# typed: strict

class Web::UsersController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :redirect_if_authenticated

  sig { returns(T.untyped) }
  def new
    user = User.new(sign_up_form_data)

    seo_tags = {
      title: t(".title"),
      canonical: view_context.new_user_url,
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      demo: params[:demo],
      user: UserSignUpFormResource.new(user)
    }
  end

  sig { returns(T.untyped) }
  def create
    struct = ApplicationParamsStruct.from_params(SignUpStruct, params.require(:data))

    return fail_sign_up(struct.first_name, struct.email, struct.errors) unless struct.valid?

    result = UserService.sign_up(struct, locale: I18n.locale.to_s)

    case result
    when Typed::Failure
      user = result.error
      return fail_sign_up(user.first_name, user.email, user.errors)
    end

    user = result.payload

    signed_up_event_data = {
      user_id: user.id,
      email: user.email,
      first_name: user.first_name,
      last_name: user.last_name,
      locale: I18n.locale.to_s
    }
    user_signed_up = UserSignedUpEvent.new(data: signed_up_event_data)

    publish_event(user_signed_up, user)
    js_event(user_signed_up)

    sign_in(user)

    lesson_ids = session.fetch(:finished_as_guest, {}).keys
    lesson_ids.each do |id|
      lesson = Language::Lesson.find(id)
      language = lesson.language
      language_member = language.members.find_or_initialize_by(user:)
      if language_member.new_record?
        language_member.save!

        course_started_event_data = {
          occurrence_count: 1,
          slug: language.slug,
          locale: I18n.locale.to_s
        }
        course_started_event = CourseStartedEvent.new(data: course_started_event_data)
        publish_event(course_started_event, user)
        js_event(course_started_event)
      end

      lesson_member = language_member.lesson_members.build(
        language:,
        lesson:,
        user:
      )
      lesson_member.finish!

      lesson_started_and_finished_event_data = {
        occurrence_count: 1,
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale: I18n.locale.to_s
      }

      started_event = LessonStartedEvent.new(data: lesson_started_and_finished_event_data)

      publish_event(started_event, user)
      js_event(started_event)

      finished_event = LessonFinishedEvent.new(data: lesson_started_and_finished_event_data)

      publish_event(finished_event, user)
      js_event(finished_event)
    end

    f(:success)
    redirect_to params[:from].presence || root_path
  end

  private

  sig { params(first_name: T.untyped, email: T.untyped, errors: T.untyped).returns(T.untyped) }
  def fail_sign_up(first_name, email, errors)
    f(:error)
    flash.inertia[:sign_up_form] = { first_name:, email: }
    redirect_to new_user_url, inertia: { errors: }
  end

  sig { returns(T.untyped) }
  def sign_up_form_data
    flash.inertia[:sign_up_form].presence || {}
  end
end
