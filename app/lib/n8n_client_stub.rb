class N8nClientStub
  extend T::Sig

  sig do
    params(
      serializer: WorkflowLeadSerializer
    ).returns(T::Hash[Symbol, T.untyped])
  end
  def trigger_lead_created_workflow(serializer)
    {}
  end
end
