# frozen_string_literal: true

class LessonMemberMutator
  def self.find_or_create_member!(lesson:, lesson_version:, language:, user:)
    member = lesson.members.find_by(language: language, user: user)
    return member unless member.nil?

    member = lesson.members.build(
      lesson_version: lesson_version,
      user: user,
      language: language
    )
    member.save!

    member
  end
end
