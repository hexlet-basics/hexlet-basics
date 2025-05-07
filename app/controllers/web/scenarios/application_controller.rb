class Web::Scenarios::ApplicationController < Web::ApplicationController
  before_action :authenticate_user!

  def resource_scenario
    @resource_scenario ||= Survey::Scenario.find(params[:scenario_id])
  end
end
