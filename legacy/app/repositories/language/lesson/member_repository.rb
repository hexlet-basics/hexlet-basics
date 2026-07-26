# typed: strict
# frozen_string_literal: true

module Language::Lesson::MemberRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { ActiveRecord::Base }

  included do
    T.bind(self, T.class_of(Language::Lesson::Member))
    scope :started_or_nil, -> { started.or(where(state: nil)) }
  end
end
