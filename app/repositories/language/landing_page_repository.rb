# typed: strict

module Language::LandingPageRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { ActiveRecord::Base }
  include LocaleRepository

  included do
    T.bind(self, T.class_of(Language::LandingPage))
    scope :web, ->() do
      published.with_locale
        .joins(:language)
        .merge(Language.completed_progress)
        .includes({ language: [ :current_version, { cover_attachment: :blob } ] })
    end
  end
end
