# frozen_string_literal: true

class Web::Admin::Languages::VersionsController < Web::Admin::Languages::ApplicationController
  def index
    @versions = resource_language.versions.order(created_at: :desc).page(params[:page])
  end

  def create
    @version = resource_language.versions.build

    if @version.save
      ExerciseLoaderJob.perform_later(@version.id)

      redirect_to admin_language_versions_path(resource_language)
    else
      render :index
    end
  end
end
