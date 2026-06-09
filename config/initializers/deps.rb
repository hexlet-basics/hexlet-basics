Rails.configuration.to_prepare do
  next if HexletBasics::Application.assets_precompiling?

  Rails.application.config.x.dependencies = DependencyFactory.build
end
