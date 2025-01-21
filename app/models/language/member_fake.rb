# frozen_string_literal: true

class Language::MemberFake
  def finished?
    false
  end

  def state
    "ready_to_start"
  end
end
