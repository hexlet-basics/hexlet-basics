# frozen_string_literal: true

class Lesson::MemberService
  def self.ensure_member!(lesson, user)
    return Language::Lesson::FakeMember.new if user.guest?

    LessonMemberMutator.find_or_create_member!(
      lesson: lesson,
      lesson_version: lesson.current_version,
      user: user,
      language: lesson.language
    )
  end
end
