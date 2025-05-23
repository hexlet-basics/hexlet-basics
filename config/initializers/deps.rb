Rails.configuration.to_prepare do
  Rails.application.config.x.dependencies = DependencyFactory.build
end
