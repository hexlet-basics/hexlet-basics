# typed: strict
# frozen_string_literal: true

# RubyLLM mixes `RubyLLM::ActiveRecord::ActsAs` into `ActiveRecord::Base`, which
# exposes these class-level macros (see the gem's `active_record/acts_as.rb`).
# Tapioca does not capture the runtime `ActiveRecord::Base.include`, so declare
# the macros here to let `acts_as_*`-based models run under `typed: strict`.
class ActiveRecord::Base
  class << self
    sig { params(kwargs: T.untyped).void }
    def acts_as_chat(**kwargs); end

    sig { params(kwargs: T.untyped).void }
    def acts_as_model(**kwargs); end

    sig { params(kwargs: T.untyped).void }
    def acts_as_message(**kwargs); end

    sig { params(kwargs: T.untyped).void }
    def acts_as_tool_call(**kwargs); end
  end
end
