# typed: strict
# frozen_string_literal: true

# Owns the course/lesson progress lifecycle for a user: starting a course,
# starting a lesson, and recording a solution check (which may finish the
# lesson and, in turn, the course). It drives the AASM transitions on
# Language::Member / Language::Lesson::Member, publishes the resulting domain
# events to RES, and returns those events so callers can selectively forward
# them to the frontend via js_event.
class CourseProgressService < ApplicationService
  class StartLessonPayload < T::Struct
    const :lesson_member, Language::Lesson::Member
    const :events, T::Array[ApplicationEvent], default: []
  end

  # The result of a solution check, shaped for LessonCheckingResponseResource:
  # the exercise outcome plus whether this check finished the lesson / course.
  class CheckPayload < T::Struct
    const :passed, T::Boolean
    const :output, String
    const :result, String
    const :status, Integer
    const :lesson_has_been_finished, T::Boolean
    const :language_has_been_finished, T::Boolean
  end

  class << self
    extend T::Sig

    # Starts the course for the user (idempotently) and then the lesson within
    # it. Course enrollment is an internal step: callers ask only to start a
    # lesson. Publishes CourseStartedEvent / LessonStartedEvent for the steps
    # that actually transitioned and returns them for selective js_event.
    sig { params(user: User, language: Language, lesson: Language::Lesson, locale: String).returns(Typed::Result[StartLessonPayload, NilClass]) }
    def start_lesson(user:, language:, lesson:, locale:)
      events = T.let([], T::Array[ApplicationEvent])

      course_member = start_course!(user, language, locale, events)
      lesson_member = start_lesson!(course_member, lesson, user, locale, events)

      success_with(StartLessonPayload.new(lesson_member:, events:))
    end

    # Runs the lesson's exercise check and records its outcome. Always publishes
    # SolutionCheckedEvent (guest-safe). On a passing check by a signed-in user,
    # finishes the lesson and, if it was the last one, the course. Returns a
    # typed payload with the exercise outcome and the finished flags.
    sig { params(user: T.nilable(User), lesson: Language::Lesson, lesson_version: T.untyped, language_version: T.untyped, code: T.untyped, locale: String).returns(Typed::Result[CheckPayload, NilClass]) }
    def record_check(user:, lesson:, lesson_version:, language_version:, code:, locale:)
      language = lesson.language
      exercise = LessonTester.run(lesson_version, language_version, code, user)

      solution_checked_event = SolutionCheckedEvent.new(data: {
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale:,
        passed: exercise.passed
      })
      EventSender.publish_event(solution_checked_event, user)

      lesson_finished, course_finished =
        if exercise.passed && user
          ActiveRecord::Base.transaction do
            course_member = language.members.find_by!(user:)
            lesson_member = lesson.members.find_by!(user:)

            [ finish_lesson!(course_member, lesson_member, lesson, language, user, locale),
             finish_course!(course_member, language, user, locale) ]
          end
        else
          [ false, false ]
        end

      success_with(CheckPayload.new(
        passed: exercise.passed,
        output: exercise.output,
        result: exercise.result,
        status: exercise.status,
        lesson_has_been_finished: lesson_finished,
        language_has_been_finished: course_finished
      ))
    end

    private

    sig { params(user: User, language: Language, locale: String, events: T::Array[ApplicationEvent]).returns(Language::Member) }
    def start_course!(user, language, locale, events)
      course_member = language.members.find_or_initialize_by(user:)
      return course_member unless course_member.new_record?

      course_member.save!

      event = CourseStartedEvent.new(data: {
        occurrence_count: user.language_members.started.count,
        slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)
      events << event

      course_member
    end

    sig { params(course_member: Language::Member, lesson: Language::Lesson, user: User, locale: String, events: T::Array[ApplicationEvent]).returns(Language::Lesson::Member) }
    def start_lesson!(course_member, lesson, user, locale, events)
      language = course_member.language
      lesson_member = course_member.lesson_members.find_or_create_by!(language:, lesson:, user:)
      return lesson_member unless lesson_member.previously_new_record?

      event = LessonStartedEvent.new(data: {
        occurrence_count: course_member.lesson_members.count,
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)
      events << event

      lesson_member
    end

    # Finishes the lesson if the AASM guard allows it. Returns whether it
    # transitioned, publishing LessonFinishedEvent when it does.
    sig { params(course_member: Language::Member, lesson_member: Language::Lesson::Member, lesson: Language::Lesson, language: Language, user: User, locale: String).returns(T::Boolean) }
    def finish_lesson!(course_member, lesson_member, lesson, language, user, locale)
      return false unless lesson_member.may_finish?

      lesson_member.finish!

      event = LessonFinishedEvent.new(data: {
        occurrence_count: course_member.lesson_members.count,
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)

      true
    end

    # Finishes the course if every lesson is done (AASM guard). Returns whether
    # it transitioned, publishing CourseFinishedEvent when it does.
    sig { params(course_member: Language::Member, language: Language, user: User, locale: String).returns(T::Boolean) }
    def finish_course!(course_member, language, user, locale)
      return false unless course_member.may_finish?

      course_member.finish!

      event = CourseFinishedEvent.new(data: {
        occurrence_count: user.language_members.finished.count,
        slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)

      true
    end
  end
end
