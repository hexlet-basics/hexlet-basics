# frozen_string_literal: true

namespace :js_routes do
  desc 'Generate Routes for JS'
  task generate: :environment do
    JsRoutes.generate!(Rails.root.join('app/packs/vendor/appRoutes.js'), camel_case: true)
  end
end
