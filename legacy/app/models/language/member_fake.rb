# typed: strict
# frozen_string_literal: true

class Language::MemberFake
  extend T::Sig

  sig { returns(T::Boolean) }
  def finished?
    false
  end

  sig { returns(T.untyped) }
  def state
    "ready_to_start"
  end
end
