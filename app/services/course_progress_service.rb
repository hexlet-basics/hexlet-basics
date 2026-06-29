# typed: strict
# frozen_string_literal: true

# Owns the course/lesson progress lifecycle for a user: starting a course,
# starting a lesson, and recording a solution check (which may finish the
# lesson and, in turn, the course). It drives the AASM transitions on
# Language::Member / Language::Lesson::Member, publishes the resulting domain
# events to RES, and returns those events so callers can selectively forward
# them to the frontend via js_event.
class CourseProgressService < ApplicationService
  class StartCoursePayload < T::Struct
    const :course_member, Language::Member
    const :events, T::Array[ApplicationEvent], default: []
  end

  class StartLessonPayload < T::Struct
    const :lesson_member, Language::Lesson::Member
    const :events, T::Array[ApplicationEvent], default: []
  end

  class CheckPayload < T::Struct
    const :events, T::Array[ApplicationEvent], default: []
  end

  class << self
    extend T::Sig

    sig { params(user: User, language: Language, locale: Symbol).returns(StartCoursePayload) }
    def start_course(user, language, locale:)
      course_member = language.members.find_or_initialize_by(user:)
      return StartCoursePayload.new(course_member:) unless course_member.new_record?

      course_member.save!

      event = CourseStartedEvent.new(data: {
        occurrence_count: user.language_members.started.count,
        slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)

      StartCoursePayload.new(course_member:, events: [ event ])
    end

    sig { params(course_member: Language::Member, lesson: Language::Lesson, user: User, locale: Symbol).returns(StartLessonPayload) }
    def start_lesson(course_member, lesson, user, locale:)
      language = course_member.language
      lesson_member = course_member.lesson_members.find_or_create_by!(language:, lesson:, user:)
      return StartLessonPayload.new(lesson_member:) unless lesson_member.previously_new_record?

      event = LessonStartedEvent.new(data: {
        occurrence_count: course_member.lesson_members.count,
        lesson_slug: lesson.slug,
        course_slug: language.slug,
        locale:
      })
      EventSender.publish_event(event, user)

      StartLessonPayload.new(lesson_member:, events: [ event ])
    end

    # Always publishes SolutionCheckedEvent (guest-safe). On a passing check by a
    # signed-in user, finishes the lesson and, if it was the last one, the course.
    sig { params(user: T.nilable(User), lesson: Language::Lesson, language: Language, passed: T::Boolean, locale: Symbol).returns(CheckPayload) }
    def record_check(user:, lesson:, language:, passed:, locale:)
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

      CheckPayload.new(events:)
    end

    private

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
