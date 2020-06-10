# frozen_string_literal: true

class Web::Languages::Modules::ApplicationController < Web::Languages::ApplicationController
  helper_method :resource_module

  def resource_module
    @resource_module ||= resource_language.modules.find_by!(slug: params[:module_id])
  end
end
