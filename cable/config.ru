# https://guides.rubyonrails.org/action_cable_overview.html#standalone
require_relative "../config/environment"
Rails.application.eager_load!

# NOTE: for kubernetes probes
map "/up" do
  run ->(_env) { [ 200, {}, [ "it works!" ] ] }
end

run ActionCable.server
