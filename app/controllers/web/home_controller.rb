# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languges_published = Language.published.includes(:current_version)
    @languages_in_development = Language.in_development.includes(:current_version)
    @language_members_by_language = current_user.language_members.index_by(&:language_id)
  end
end
