class Web::Admin::ApplicationController < Web::ApplicationController
  before_action :authenticate_admin!
end
