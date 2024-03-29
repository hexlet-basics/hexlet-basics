# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `slim-rails` gem.
# Please instead update this file by running `bin/tapioca gem slim-rails`.

module Slim; end
class Slim::InvalidAttributeNameError < ::StandardError; end
module Slim::Rails; end
class Slim::Rails::Railtie < ::Rails::Railtie; end

module Slim::Rails::RegisterEngine
  class << self
    # source://slim-rails//lib/slim-rails/register_engine.rb#12
    def register_engine(app, config); end

    private

    # source://slim-rails//lib/slim-rails/register_engine.rb#29
    def _register_engine(config); end

    # source://slim-rails//lib/slim-rails/register_engine.rb#22
    def _register_engine3(app); end
  end
end

class Slim::Rails::RegisterEngine::Transformer
  class << self
    # source://slim-rails//lib/slim-rails/register_engine.rb#5
    def call(input); end
  end
end

# source://slim-rails//lib/slim-rails/version.rb#3
Slim::Rails::VERSION = T.let(T.unsafe(nil), String)

class Slim::RailsTemplate < ::Temple::Templates::Rails; end
class Slim::Template < ::Temple::Templates::Tilt; end

# source://slim/4.1.0/lib/slim/version.rb#4
Slim::VERSION = T.let(T.unsafe(nil), String)
