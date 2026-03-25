class Web::Account::ApplicationController < Web::ApplicationController
  before_action :require_authentication
end
