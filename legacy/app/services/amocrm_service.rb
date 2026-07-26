# typed: strict

class AmocrmService
  extend T::Sig

  sig { params(lead: Lead).void }
  def self.create_lead(lead:)
    payload = WorkflowLeadResource.new(lead).to_h
    DepsLocator.current.amocrm.unsorted_leads.create_forms(body: payload["lead"])
  end
end
