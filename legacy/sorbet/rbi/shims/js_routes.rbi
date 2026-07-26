# typed: true

# js-routes is excluded from tapioca gem RBI generation (its auto-RBI trips
# Sorbet 5070 on `T.nilable(T.untyped)`, and because the gem is excluded
# tapioca deletes any hand-maintained `sorbet/rbi/gems/js-routes@*.rbi` as an
# orphan on every `tapioca gem` run). We only use `JsRoutes.setup` (in
# config/initializers/js_routes.rb), so shim just that — shims live outside
# sorbet/rbi/gems and survive regeneration.
module JsRoutes
  class << self
    sig { params(block: T.untyped).void }
    def setup(&block); end
  end
end
