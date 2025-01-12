class ReviewResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Review

  # one :user, resource: UserResource
  has_one :user, serializer: UserResource

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id, :created_at, :locale
end
