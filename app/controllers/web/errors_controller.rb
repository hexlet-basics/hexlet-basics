# frozen_string_literal: true

class Web::ErrorsController < Web::ApplicationController
  before_action do
    set_meta_tags title: t('.base')
  end

  def forbidden
    respond_to do |format|
      format.html { render status: :forbidden }
      format.all  { render nothing: true, status: :forbidden }
    end
  end

  def not_found
    respond_to do |format|
      format.html { render status: :not_found }
      format.all  { head :not_found }
    end
  end

  def server_error
    respond_to do |format|
      format.html { render status: :internal_server_error }
      format.all  { render nothing: true, status: :internal_server_error }
    end
  end
end
