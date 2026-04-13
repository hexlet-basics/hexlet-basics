# typed: strict

module DepsLocator
  extend T::Sig

  sig { returns(Dependencies) }
  def self.current
    dependencies = Rails.application.config.x.dependencies
    return dependencies if dependencies.is_a?(Dependencies)

    DependencyFactory.build.tap do |built_dependencies|
      Rails.application.config.x.dependencies = built_dependencies
    end
  end
end
