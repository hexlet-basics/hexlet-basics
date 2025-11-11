class ApplicationController < ActionController::Base
  include ActiveStorage::SetCurrent
  include Pundit::Authorization
  include Pagy::Method
  # rescue_from Pagy::OverflowError, with: :redirect_to_last_page
  # rescue_from Pagy::VariableError, with: :redirect_to_last_page

  include AuthConcern
  include AhoyConcern

  def default_url_options
    { suffix: params[:suffix] }
  end

  # def redirect_to_last_page(exception)
  #   redirect_to url_for(page: exception.pagy.last), notice: "Page ##{params[:page]} is overflowing. Showing page #{exception.pagy.last} instead."
  # end
end
