# frozen_string_literal: true

class Web::Admin::Languages::VersionsController < Web::Admin::Languages::ApplicationController
  def index
    @q = resource_language.versions.ransack(params[:q])
    @q.sorts = 'created_at desc' if @q.sorts.empty?
    @versions = @q.result.page(params[:page])
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
