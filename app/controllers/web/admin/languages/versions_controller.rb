# typed: strict
# frozen_string_literal: true

class Web::Admin::Languages::VersionsController < Web::Admin::Languages::ApplicationController
  # def index
  #   q = params.fetch(:q, {}).with_defaults("s" => "created_at desc")
  #   @search = resource_language.versions.ransack(q)
  #   @versions = @search.result.page(params[:page])
  #
  #   render inertia: true, props: {}
  # end

  sig { returns(T.untyped) }
  def create
    @version = T.let(resource_language.versions.build, T.untyped)

    if @version.save
      ExerciseLoaderJob.perform_later(@version.id)
      f(:success)
    else
      f(:error)
    end

    redirect_to edit_admin_language_path(resource_language)
  end
end
