# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languages_completed = Language.with_progress(:completed)
                                   .joins(current_version: :infos)
                                   .where(current_version: { language_version_infos: { locale: I18n.locale } })
    @languages_in_development = Language.with_progress(:in_development).includes(:current_version)
    @language_members_by_language = current_user.language_members.index_by(&:language_id)
    # @language_member_fake = Language::MemberFake.new
  end
end
