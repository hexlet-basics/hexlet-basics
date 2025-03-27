module Language::LandingPageRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
    scope :web, ->() do
      with_locale
        .joins(:language)
        .merge(Language.completed_progress)
        .includes({ language: [ :current_version, { cover_attachment: :blob } ] })
    end
  end
end
