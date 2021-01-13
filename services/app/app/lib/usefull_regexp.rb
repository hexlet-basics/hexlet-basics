# frozen_string_literal: true

class UsefullRegexp
  class << self
    def alphanumeric_for_all_langs
      /\A[\p{L}\w\s.-]+\z/u
    end
  end
end
