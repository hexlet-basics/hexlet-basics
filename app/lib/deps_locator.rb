# typed: strict

module DepsLocator
  extend T::Sig

  sig { returns(Dependencies) }
  def self.current
    Rails.application.config.x.dependencies
  end
end
