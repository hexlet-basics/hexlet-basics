# frozen_string_literal: true

class Api::LanguagesController < Api::ApplicationController
  def index
    @languages = Language.with_progress(:completed)
                         .joins(current_version: :infos)
                         .merge(Language::Version::Info.with_locale)

    @infos_by_language = Language::Version::Info.where(language_version: @languages.pluck(:current_version_id))
                                                .with_locale
                                                .index_by(&:language_id)
  end

  def show
    @language = Language.includes(:current_version)
                        .with_progress(:completed)
                        .find(params[:id])

    @language_info = @language.current_version.infos.find_by!(locale: I18n.locale)
  end
end
