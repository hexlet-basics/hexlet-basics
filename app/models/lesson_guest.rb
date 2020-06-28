# frozen_string_literal: true

class LessonGuest
  attr_reader :lesson

  def initialize(lesson)
    @lesson = lesson
  end

  def lesson_version
    lesson.current_version
  end
end
