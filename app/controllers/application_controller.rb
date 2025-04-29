class ApplicationController < ActionController::Base
  include ActiveStorage::SetCurrent
  include Pundit::Authorization
  include Pagy::Backend

  include AuthConcern

  def default_url_options
    { suffix: params[:suffix] }
  end
end
