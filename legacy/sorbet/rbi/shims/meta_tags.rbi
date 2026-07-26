# typed: true

# The meta-tags gem mixes `MetaTags::ViewHelper` into `ActionView::Base` through
# its Railtie `on_load(:action_view)` hook, which tapioca can't observe. State
# the include here so view helpers (`display_meta_tags`, `set_meta_tags`, ...)
# resolve on `helpers` (an `ActionView::Base` subclass) without `T.unsafe`.
class ActionView::Base
  include ::MetaTags::ViewHelper
end
