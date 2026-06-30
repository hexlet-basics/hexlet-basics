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

  class CheckPayload < T::Struct
    const :exercise_data, T::Hash[Symbol, T.untyped]
    const :events, T::Array[ApplicationEvent], default: []
  end

  class << self
    extend T::Sig

    # Starts the course for the user (idempotently) and then the lesson within
    # it. Course enrollment is an internal step: callers ask only to start a
    # lesson. Publishes CourseStartedEvent / LessonStartedEvent for the steps
    # that actually transitioned and returns them for selective js_event.
    sig { params(user: User, language: Language, lesson: Language::Lesson, locale: Symbol).returns(StartLessonPayload) }
    def start_lesson(user:, language:, lesson:, locale:)
      events = T.let([], T::Array[ApplicationEvent])

      course_member = start_course!(user, language, locale, events)
      lesson_member = start_lesson!(course_member, lesson, user, locale, events)

      StartLessonPayload.new(lesson_member:, events:)
    end

    # Runs the lesson's exercise check and records its outcome. Always publishes
    # SolutionCheckedEvent (guest-safe). On a passing check by a signed-in user,
    # finishes the lesson and, if it was the last one, the course. Returns the
    # exercise data alongside the produced events.
    sig { params(user: T.nilable(User), lesson: Language::Lesson, lesson_version: T.untyped, language_version: T.untyped, code: T.untyped, locale: Symbol).returns(CheckPayload) }
    def record_check(user:, lesson:, lesson_version:, language_version:, code:, locale:)
      language = lesson.language
      exercise_data = LessonTester.run(lesson_version, language_version, code, user)
      passed = exercise_data[:passed]

      events = T.let([], T::Array[ApplicationEvent])

      solution_checked_event = SolutionCheckedEvent.new(data: {
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale:,
        passed:
      })
      EventSender.publish_event(solution_checked_event, user)
      events << solution_checked_event

      if passed && user
        ActiveRecord::Base.transaction do
          course_member = language.members.find_by!(user:)
          lesson_member = lesson.members.find_by!(user:)

          lesson_finished_event = process_lesson!(course_member, lesson_member, lesson, language, user, locale)
          events << lesson_finished_event if lesson_finished_event

          course_finished_event = process_course!(course_member, language, user, locale)
          events << course_finished_event if course_finished_event
        end
      end

      CheckPayload.new(exercise_data:, events:)
    end

    private

    sig { params(user: User, language: Language, locale: Symbol, events: T::Array[ApplicationEvent]).returns(Language::Member) }
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

    sig { params(course_member: Language::Member, lesson: Language::Lesson, user: User, locale: Symbol, events: T::Array[ApplicationEvent]).returns(Language::Lesson::Member) }
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

    sig { params(course_member: Language::Member, lesson_member: Language::Lesson::Member, lesson: Language::Lesson, language: Language, user: User, locale: Symbol).returns(T.nilable(LessonFinishedEvent)) }
    def process_lesson!(course_member, lesson_member, lesson, language, user, locale)
      return unless lesson_member.may_finish?

      lesson_member.finish!

      event = LessonFinishedEvent.new(data: {
        occurrence_count: course_member.lesson_members.count,
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)

      event
    end

    sig { params(course_member: Language::Member, language: Language, user: User, locale: Symbol).returns(T.nilable(CourseFinishedEvent)) }
    def process_course!(course_member, language, user, locale)
      return unless course_member.may_finish?

      course_member.finish!

      event = CourseFinishedEvent.new(data: {
        occurrence_count: user.language_members.finished.count,
        slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)

      event
    end
  end
end
