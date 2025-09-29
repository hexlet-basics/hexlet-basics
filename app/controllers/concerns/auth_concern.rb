# frozen_string_literal: true

module AuthConcern
  def sign_up(user)
    signed_up_event_data = {
      id: user.id,
      email: user.email,
      first_name: user.first_name,
      last_name: user.last_name
    }
    user_signed_up = UserSignedUpEvent.new(data: signed_up_event_data)

    publish_event(user_signed_up, user)
    event_to_js(user_signed_up)
  end

  def sign_in(user)
    session[:user_id] = user.id
    ahoy.authenticate(user)

    event_data = {
      id: user.id,
      email: user.email,
      first_name: user.first_name,
      last_name: user.last_name
    }
    event = UserSignedInEvent.new(data: event_data)
    publish_event(event, user)
  end

  def fill_guests_data(user)
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
          locale: I18n.locale
        }
        course_started_event = CourseStartedEvent.new(data: course_started_event_data)
        publish_event(course_started_event, user)
        event_to_js(course_started_event)
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
        locale: I18n.locale
      }

      started_event = LessonStartedEvent.new(data: lesson_started_and_finished_event_data)

      publish_event(started_event, user)
      event_to_js(started_event)

      finished_event = LessonFinishedEvent.new(data: lesson_started_and_finished_event_data)

      publish_event(finished_event, user)
      event_to_js(finished_event)
    end
  end

  def sign_out
    session.delete(:user_id)
    session.clear
  end

  def signed_in?
    !current_user.guest?
  end

  def current_user
    id = session[:user_id]
    @current_user ||= id && User.active.find_by(id: id) || Guest.new
  end

  def authenticate_user!
    return if signed_in?

    redirect_to new_session_path
  end

  def guests_only!
    return if !signed_in?

    redirect_to root_path
  end

  def authenticate_admin!
    redirect_to root_path unless current_user.admin?
  end

  def require_api_auth!
    head :forbidden unless signed_in?
  end

  def require_admin_api_auth!
    head :forbidden unless current_user.admin?
  end
end
