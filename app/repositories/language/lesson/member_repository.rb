# frozen_string_literal: true

module Language::Lesson::MemberRepository
  extend ActiveSupport::Concern

  included do
    scope :started_or_nil, -> { started.or(where(state: nil)) }
  end
end
