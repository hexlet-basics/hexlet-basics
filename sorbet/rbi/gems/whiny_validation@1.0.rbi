# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `whiny_validation` gem.
# Please instead update this file by running `bin/tapioca gem whiny_validation`.

module WhinyValidation
  extend ::ActiveSupport::Concern

  # source://whiny_validation//lib/whiny_validation.rb#14
  def whiny_validation; end

  class << self
    # source://whiny_validation//lib/whiny_validation/configuration.rb#2
    def configuration; end

    # @yield [configuration]
    #
    # source://whiny_validation//lib/whiny_validation/configuration.rb#6
    def configure; end
  end
end

class WhinyValidation::Configuration
  # @return [Configuration] a new instance of Configuration
  #
  # source://whiny_validation//lib/whiny_validation/configuration.rb#13
  def initialize; end

  # Returns the value of attribute log_level.
  #
  # source://whiny_validation//lib/whiny_validation/configuration.rb#11
  def log_level; end

  # Sets the attribute log_level
  #
  # @param value the value to set the attribute log_level to.
  #
  # source://whiny_validation//lib/whiny_validation/configuration.rb#11
  def log_level=(_arg0); end
end

class WhinyValidation::LogSubscriber < ::ActiveSupport::LogSubscriber
  # source://whiny_validation//lib/whiny_validation.rb#21
  def validation_failed(event); end
end

# source://whiny_validation//lib/whiny_validation/version.rb#2
WhinyValidation::VERSION = T.let(T.unsafe(nil), String)