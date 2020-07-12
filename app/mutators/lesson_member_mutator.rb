# frozen_string_literal: true

class LessonMemberMutator
  def self.find_or_create_member!(lesson:, lesson_version:, language:, user:)
    member = lesson.find_by(language: language, user: user)
    return member unless member.nil?

    member = lesson.members.build(
      lesson: lesson,
      lesson_version: lesson_version,
      user: user,
      language: language
    )
    member.save!

    member
  end
end
