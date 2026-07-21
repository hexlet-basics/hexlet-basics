# typed: true

# `identified_by :current_user` in ApplicationCable::Connection defines the
# `current_user` reader/writer dynamically. Tapioca has no DSL compiler for
# `identified_by`, so shim the accessors we rely on to keep the connection at
# `# typed: strict`.
module ApplicationCable
  class Connection < ActionCable::Connection::Base
    sig { returns(T.nilable(::User)) }
    def current_user; end

    sig { params(value: T.nilable(::User)).returns(T.nilable(::User)) }
    def current_user=(value); end
  end
end
