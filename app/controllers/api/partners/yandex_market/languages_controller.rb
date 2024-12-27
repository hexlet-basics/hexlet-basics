# frozen_string_literal: true

class Api::Partners::YandexMarket::LanguagesController < Api::Partners::YandexMarket::ApplicationController
  def index
    @languages = Language.joins(current_version: :infos)
                         .includes([ :current_version ])
                         .merge(Language::Version::Info.with_locale)

    @infos_by_language = Language::Version::Info.where(language_version: @languages.pluck(:current_version_id))
                                                .with_locale
                                                .index_by(&:language_id)
  end
end
