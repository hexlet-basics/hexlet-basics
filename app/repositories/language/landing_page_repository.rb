module Language::LandingPageRepository
  extend ActiveSupport::Concern
  include LocaleRepository

  included do
    scope :web, ->() do
      published.with_locale
        .joins(:language)
        .merge(Language.completed_progress)
        .includes({ language: [ :current_version, { cover_attachment: :blob } ] })
    end
  end
end
