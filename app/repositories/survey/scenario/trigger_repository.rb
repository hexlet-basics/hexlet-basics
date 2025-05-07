module Survey::Scenario::TriggerRepository
  extend ActiveSupport::Concern

  included do
    scope :with_event_threshold_met, ->(occurrence_count) {
      where("event_threshold_count IS NULL OR event_threshold_count::int <= ?", occurrence_count)
    }
  end
end
