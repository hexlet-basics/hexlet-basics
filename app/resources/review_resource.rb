class ReviewResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Review

  # one :user, resource: UserResource
  has_one :user
  has_one :language, resource: LanguageResource

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id, :created_at, :locale

  typelize :string, nullable: true
  attribute :full_name do
    [ it.first_name, it.last_name ].join " "
  end

  # typelize :Array, nullable: true
  # attribute :state_events do |obj|
  #   obj.class.enum_as_hashes(:states)
  # end
end
