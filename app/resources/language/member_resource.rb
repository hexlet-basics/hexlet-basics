class Language::MemberResource
  # include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Member

  attributes :id, :state, :language_id
  # has_one :current_version, resource: Language::VersionResource

  # typelize :string, nullable: true
  # attribute :repository_url do |obj|
  #   obj.repository_url
  # end
end
