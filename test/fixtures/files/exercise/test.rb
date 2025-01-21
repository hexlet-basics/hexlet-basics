# frozen_string_literal: true

require "minitest-power_assert"

require "minitest/pride"
require "minitest/autorun"

require ARGV[0]

describe "function" do
  it "should work" do
    assert { sum(2, 2) == 4 }
  end
end
