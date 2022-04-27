# frozen_string_literal: true

class Api::LanguagesController < Api::ApplicationController
  def index
    @languages = Language.with_progress(:completed)
                         .joins(current_version: :infos)
                         .where(current_version: { language_version_infos: { locale: I18n.locale } })

    @infos_by_language = Language::Version::Info.where(language_version: @languages.pluck(:current_version_id))
                                                .with_locale
                                                .index_by(&:language_id)
  end

  def show
    @language = Language.includes(:current_version)
                        .with_progress(:completed)
                        .find(params[:id])

    @language_info = @language.current_version.infos.with_locale.first
  end
end
