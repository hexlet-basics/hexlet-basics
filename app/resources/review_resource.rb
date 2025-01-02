class ReviewResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :id, :body, :created_at

  one :language, resource: LanguageResource

  typelize :string, nullable: true
  attribute :full_name do |review|
    [ review.first_name, review.last_name ].join(" ")
  end
end
