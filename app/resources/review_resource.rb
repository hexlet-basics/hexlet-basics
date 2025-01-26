class ReviewResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Review

  # one :user, resource: UserResource
  has_one :user
  has_one :language, resource: OriginalLanguageResource

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id, :created_at, :locale

  typelize :string, nullable: true
  attribute :full_name do |review|
    [ review.first_name, review.last_name ].join " "
  end

  typelize "Array<[string, string]>"
  attribute :state_events do |obj|
    obj.aasm.events_for_select
  end
end
