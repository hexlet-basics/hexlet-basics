# typed: true

# minitest-power_assert overrides Minitest::Assertions#assert to accept a block
# (`assert { expression }`). Tapioca generates that block form on
# Minitest::PowerAssert::Assertions, but tests resolve `assert` against the base
# Minitest::Assertions, whose RBI is `assert(test, msg = nil)` (no block). Shim the
# base to accept an optional value + block so `assert { ... }` typechecks.
module Minitest::Assertions
  sig { params(test: T.untyped, msg: T.untyped, blk: T.untyped).returns(T.untyped) }
  def assert(test = T.unsafe(nil), msg = T.unsafe(nil), &blk); end
end
