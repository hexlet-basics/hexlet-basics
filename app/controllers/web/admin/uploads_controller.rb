# frozen_string_literal: true

class Web::Admin::UploadsController < Web::Admin::ApplicationController
  def index
    @uploads = Upload.all
  end

  def new
    @upload = Upload.new
  end

  def create
    @upload = Upload.new(upload_params)

    if @upload.save
      ExerciseLoaderWorker.perform_async(@upload.id)

      redirect_to admin_uploads_path
    else
      render :new
    end
  end

  private

  def upload_params
    params.require(:upload).permit(:language_name)
  end
end
