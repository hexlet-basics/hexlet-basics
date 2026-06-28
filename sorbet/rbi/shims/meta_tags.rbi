# typed: true

# The meta-tags gem mixes MetaTags::ControllerHelper (set_meta_tags, meta_tags)
# into ActionController::Base at boot via its Railtie. tapioca's gem RBI defines
# the module but not the runtime include, so typed controllers don't see
# set_meta_tags. Re-declare the include here so it resolves for all controllers.
class ActionController::Base
  include ::MetaTags::ControllerHelper
end
