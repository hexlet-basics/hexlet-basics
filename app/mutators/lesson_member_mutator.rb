# frozen_string_literal: true

class LessonMemberMutator
  def self.find_or_create_member!(lesson, user)
    member = lesson.members.find_by(language: lesson.language, user: user)
    return member if !member.nil?

    member = lesson.members.build(
      lesson_version: lesson.current_version,
      user: user,
      language: lesson.language
    )
    member.save!

    member
  end
end
