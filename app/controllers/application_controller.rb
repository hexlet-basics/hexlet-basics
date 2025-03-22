class ApplicationController < ActionController::Base
  include ActiveStorage::SetCurrent
  include Pagy::Backend

  include AuthConcern

  def default_url_options
    { suffix: params[:suffix] }
  end
end
