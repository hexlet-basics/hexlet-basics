# frozen_string_literal: true

class Guest
  def guest?
    true
  end

  def admin?
    false
  end

  def finished_members_for_language(_)
    []
  end

  def complete_language?(*)
    false
  end
end
