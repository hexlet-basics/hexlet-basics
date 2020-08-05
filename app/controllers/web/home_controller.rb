# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languages = Language.includes(:current_version)
    @language_members_by_language = current_user.language_members.index_by(&:language_id)
  end
end
