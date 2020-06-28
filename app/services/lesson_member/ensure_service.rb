# frozen_string_literal: true

class LessonMember::EnsureService
  attr_reader :lesson, :user

  def initialize(lesson, user)
    @lesson = lesson
    @user = user
  end

  def execute
    return LessonGuest.new(lesson) if user.guest?

    find_lesson_member || create_lesson_member
  end

  def find_lesson_member
    LessonMember.find_by(lesson: lesson, user: user)
  end

  def create_lesson_member
    LessonMember.create!(lesson: lesson, lesson_version: lesson.current_version, user: user)
  end
end
