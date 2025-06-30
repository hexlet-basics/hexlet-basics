class Ahoy::Store < Ahoy::DatabaseStore
end

# set to true for JavaScript tracking
Ahoy.api = false
Ahoy.track_bots = Rails.env.development?

Ahoy.exclude_method = ->(_controller, request) do
  # NOTE игнорируем все запросы к любым эндпоинтам Active Storage
  request.path.start_with?("/rails/active_storage")
end

# set to true for geocoding (and add the geocoder gem to your Gemfile)
# we recommend configuring local geocoding as well
# see https://github.com/ankane/ahoy#geocoding
Ahoy.geocode = false
