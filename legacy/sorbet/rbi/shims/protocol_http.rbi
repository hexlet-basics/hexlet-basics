# typed: true

# protocol-http / async-http are transitive deps (via falcon) with no generated
# gem RBIs. webmock's auto-RBI declares `Async::HTTP::Client < Protocol::HTTP::Methods`,
# so `Protocol::HTTP::Methods` must be a Class — but `tapioca todo` otherwise guesses
# it as a module, which trips Sorbet 5067 (super class does not derive from Class).
# Shim it as a class so the webmock async adapter chain resolves; shims survive
# `tapioca todo`/`tapioca gem` regeneration.
module Protocol; end
module Protocol::HTTP; end
class Protocol::HTTP::Methods; end
