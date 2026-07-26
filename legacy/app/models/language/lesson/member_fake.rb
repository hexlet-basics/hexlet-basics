# typed: strict
# frozen_string_literal: true

class Language::Lesson::MemberFake
  extend T::Sig

  sig { returns(T.untyped) }
  def state
    "finished"
  end
end
