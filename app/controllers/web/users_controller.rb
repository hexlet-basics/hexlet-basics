# frozen_string_literal: true

class Web::UsersController < Web::ApplicationController
  def new
    user = User::SignUpForm.new

    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      demo: params[:demo],
      user: UserSignUpFormResource.new(user)
    }
  end

  def create
    user = User::SignUpForm.new(params[:user_sign_up_form])

    if user.save
      event_data = {
        id: user.id,
        email: user.email,
        first_name: user.first_name
      }
      course_started = UserSignedUpEvent.new(data: event_data)

      publish_event(course_started, user)
      event_to_js(course_started)

      sign_in user

      lesson_ids = session.fetch(:finished_as_guest, {}).keys
      lesson_ids.each do |id|
        lesson = Language::Lesson.find(id)
        language = lesson.language
        language_member = language.members.find_or_initialize_by(user:)
        if language_member.new_record?
          language_member.save!

          event_data = {
            slug: language.slug,
            locale: I18n.locale
          }
          course_started = CourseStartedEvent.new(data: event_data)
          publish_event(course_started, user)
          event_to_js(course_started)
        end

        lesson_member = language_member.lesson_members.build(
          language:,
          lesson:,
          user:
        )
        lesson_member.finish!

        event_data = {
          lesson_slug: lesson.slug,
          course_slug: language.slug,
          locale: I18n.locale
        }

        started_event = LessonStartedEvent.new(data: event_data)

        publish_event(started_event, user)
        event_to_js(started_event)

        finished_event = LessonFinishedEvent.new(data: event_data)

        publish_event(finished_event, user)
        event_to_js(finished_event)
      end

      f(:success)
      redirect_to params[:from].presence || root_path
    else
      f(:error)
      redirect_to_inertia new_user_url, user
    end
  end
end
