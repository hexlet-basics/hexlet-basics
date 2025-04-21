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
      event = UserSignedUpEvent.new(data: user.slice(:id, :email, :first_name))

      publish_event(event, user)
      event_to_js(event)

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
          event = CourseStartedEvent.new(data: event_data)
          publish_event(event, user)
          event_to_js(event)
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
        event = LessonStartedEvent.new(data: event_data)

        publish_event(event, user)
        event_to_js(event)
      end

      f(:success)
      redirect_to root_url
    else
      f(:error)
      redirect_to_inertia new_user_url, user
    end
  end
end
