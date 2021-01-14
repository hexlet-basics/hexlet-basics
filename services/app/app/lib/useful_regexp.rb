# frozen_string_literal: true

class UsefulRegexp
  class << self
    def without_spec_chars
      /\A[^`!@#$%\^&*+=]+\z/u

      # any character except:
      # '`', '!', '@', '#', '\$', '%', '\^', '&', '*', '+', '=',
    end
  end
end
