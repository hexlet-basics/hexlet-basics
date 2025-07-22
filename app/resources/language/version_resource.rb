class Language::VersionResource < ApplicationResource
  typelize_from Language::Version

  attributes :id, :result, :state, :created_at

  # typelize :id, nullable: true
  # attribute :id do |info|
  #   info.version.module.id
  # end
end
