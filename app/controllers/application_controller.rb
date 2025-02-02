class ApplicationController < ActionController::Base
  include ActiveStorage::SetCurrent
  # respond_to :html
  include Pagy::Backend

  include AuthConcern

  def default_url_options
    { suffix: params[:suffix] }
  end

  def event_store
    Rails.configuration.event_store
  end
end
