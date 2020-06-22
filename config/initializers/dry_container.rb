# frozen_string_literal: true

require 'dry/system/container'

class AppContainer < Dry::System::Container
  configure do |config|
    config.root = (Pathname.pwd + 'app')

    config.auto_register = 'lib'
  end

  load_paths!('lib')
end

Import = AppContainer.injector

AppContainer.finalize! if Rails.env.production?
