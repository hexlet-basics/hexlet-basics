class BlogPostResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id, :name, :slug, :description

  typelize :string, nullable: true
  attribute :cover do |post|
    post.cover
  end
end
