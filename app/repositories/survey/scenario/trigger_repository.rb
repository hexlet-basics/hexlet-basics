# typed: strict

module Survey::Scenario::TriggerRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { ActiveRecord::Base }

  included do
    T.bind(self, T.class_of(ActiveRecord::Base))
    scope :with_event_threshold_met, ->(occurrence_count) {
      where("event_threshold_count IS NULL OR event_threshold_count::int <= ?", occurrence_count)
    }
  end
end
