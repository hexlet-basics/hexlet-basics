# set to true for JavaScript tracking
Ahoy.api = true
Ahoy.server_side_visits = :when_needed
Ahoy.cookie_domain = :all

# set to true for geocoding (and add the geocoder gem to your Gemfile)
# we recommend configuring local geocoding as well
# see https://github.com/ankane/ahoy#geocoding
Ahoy.geocode = false

Safely.report_exception_method = ->(e) { Sentry.capture_exception(e) }
