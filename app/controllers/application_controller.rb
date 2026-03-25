class ApplicationController < ActionController::Base
  include ActiveStorage::SetCurrent
  include Pundit::Authorization
  include Pagy::Method
  include Authentication
  include LocaleConcern
  include EventConcern

  # rescue_from Pagy::OverflowError, with: :redirect_to_last_page
  # rescue_from Pagy::VariableError, with: :redirect_to_last_page

  include AuthConcern

  def default_url_options
    { suffix: params[:suffix] }
  end

  # def redirect_to_last_page(exception)
  #   redirect_to url_for(page: exception.pagy.last), notice: "Page ##{params[:page]} is overflowing. Showing page #{exception.pagy.last} instead."
  # end
end
