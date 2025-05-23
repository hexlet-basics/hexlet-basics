class LeadForm < Lead
  include ActiveFormModel

  attr_accessor :contact_method, :contact_value
  permit :contact_method, :contact_value

  validates :contact_method, presence: true
  validates :contact_value, presence: true

  after_initialize do
    if contact_method.present?
      write_attribute contact_method, contact_value
    end
  end
end
