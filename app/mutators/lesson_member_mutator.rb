# frozen_string_literal: true

class LessonMemberMutator
  def self.upsert_member!(lesson, user)
    member = lesson.members.find_by(language: lesson.language, user: user)

    if member
      member.lesson_version = lesson.current_version
      member.save!

      return member
    end

    member = lesson.members.build(
      lesson_version: lesson.current_version,
      user: user,
      language: lesson.language
    )
    member.save!

    member
  end
end
