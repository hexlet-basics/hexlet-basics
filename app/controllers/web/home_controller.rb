# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languages_completed = Language.with_progress(:completed)
                                   .joins(current_version: :infos)
                                   .merge(Language::Version::Info.with_locale)

    @languages_in_development = Language.with_progress(:in_development).includes(:current_version)
    @language_members_by_language = current_user.language_members.index_by(&:language_id)

    language_versions = Language::Version.where(current_language: @languages_completed)

    gon.languages_for_widget = language_versions.pluck(:name)
  end

  def robots
    respond_to :text
  end
end
