# frozen_string_literal: true

class Web::Admin::UploadsController < Web::Admin::ApplicationController
  def index
    @uploads = Language::Upload.includes(:language).order(created_at: :desc)
  end
end
