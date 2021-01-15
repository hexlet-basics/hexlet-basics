# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  before_action do
    title :base
  end

  def index
    @languages_completed = Language.with_progress(:completed).includes(:current_version)
    @languages_in_development = Language.with_progress(:in_development).includes(:current_version)
    @language_members_by_language = current_user.language_members.index_by(&:language_id)
  end
end
