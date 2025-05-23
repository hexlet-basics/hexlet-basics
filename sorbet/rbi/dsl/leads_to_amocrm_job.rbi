# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for dynamic methods in `LeadsToAmocrmJob`.
# Please instead update this file by running `bin/tapioca dsl LeadsToAmocrmJob`.


class LeadsToAmocrmJob
  class << self
    sig do
      params(
        payload: T.untyped,
        block: T.nilable(T.proc.params(job: LeadsToAmocrmJob).void)
      ).returns(T.any(LeadsToAmocrmJob, FalseClass))
    end
    def perform_later(payload, &block); end

    sig { params(payload: T.untyped).returns(T.untyped) }
    def perform_now(payload); end
  end
end
