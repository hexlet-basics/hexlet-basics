# typed: strict

class N8nClientInterface
  extend T::Sig
  extend T::Helpers
  abstract!

  sig do
    abstract.params(
      serializer: WorkflowLeadSerializer
    ).returns(T::Hash[Symbol, T.untyped])
  end
  def trigger_lead_created_workflow(serializer); end
end
