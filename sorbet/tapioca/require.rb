# typed: true
# frozen_string_literal: true

require "English"
require "action_cable/engine"
# ActionCable's Channel::Base / Connection::Base are autoloaded by the engine and
# otherwise absent when tapioca reflects the gem, leaving an empty actioncable RBI
# (and our channels' superclasses unresolved). Force-load them so tapioca captures
# them as real classes.
require "action_cable/channel/base"
require "action_cable/connection/base"
# Load every ActiveSupport core extension so tapioca captures methods like
# TimeWithZone#beginning_of_day that are otherwise autoloaded and missing.
require "active_support/all"
require "action_controller/railtie"
require "action_mailer/railtie"
require "action_view/railtie"
require "active_job/railtie"
require "active_model/railtie"
require "active_record/railtie"
require "active_storage/engine"
require "active_support/core_ext/integer/time"
require "bootsnap/setup"
require "bundler/setup"
require "capybara/cuprite"
require "minitest-power_assert"
require "minitest/autorun"
require "minitest/pride"
require "open3"
require "pty"
require "rails"
require "rails/test_help"
require "rails/test_unit/railtie"
require "securerandom"
require "tapioca/dsl"
require "webmock/minitest"
