# typed: true
# frozen_string_literal: true

require "English"
# Default Rails bundle (as `rails new` generates) — pulls in every framework so
# tapioca reflects complete RBIs instead of the per-railtie subset.
require "rails/all"
# ActionCable's Channel::Base / Connection::Base are autoloaded by the engine and
# otherwise absent when tapioca reflects the gem, leaving an empty actioncable RBI
# (and our channels' superclasses unresolved). Force-load them so tapioca captures
# them as real classes.
require "action_cable/channel/base"
require "action_cable/connection/base"
require "bootsnap/setup"
require "bundler/setup"
require "capybara/cuprite"
require "minitest-power_assert"
require "minitest/autorun"
require "minitest/pride"
require "open3"
require "pty"
require "rails/test_help"
require "securerandom"
require "boba"
require "tapioca/dsl"
require "webmock/minitest"
