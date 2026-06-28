# typed: strict

class BannerResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Banner

  attributes :id, :locale, :body, :url, :background, :state, :starts_at, :finishes_at, :created_at
end
