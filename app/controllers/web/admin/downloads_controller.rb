# frozen_string_literal: true

class Web::Admin::DownloadsController < Web::Admin::ApplicationController
  def create
    exercise = params[:exercise]
    DownloadJob.perform_now(exercise)
  end
end
