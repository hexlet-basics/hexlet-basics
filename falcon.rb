#!/usr/bin/env -S falcon-host
# frozen_string_literal: true

require "async/http/endpoint"
require "async/http/protocol/http2"
require "falcon/environment/rack"

hostname = File.basename(__dir__)

service hostname do
  include Falcon::Environment::Rack

  preload "preload.rb"

  # Match previous Rails/Puma behavior: single worker locally unless explicitly configured.
  count [ ENV.fetch("WEB_CONCURRENCY", 1).to_i, 1 ].max

  # Keep internal app port stable for reverse proxy/probes.
  # Use APP_SERVER_PORT for explicit overrides to avoid accidental PORT leakage from shell.
  port { ENV.fetch("APP_SERVER_PORT", 3000).to_i }

  endpoint do
    Async::HTTP::Endpoint
      .parse("http://0.0.0.0:#{port}")
      .with(protocol: Async::HTTP::Protocol::HTTP2)
  end
end
