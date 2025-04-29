class UserResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User

  attributes :id, :first_name, :email, :created_at, :last_name, :admin, :assistant_messages_count

  typelize :boolean, nullable: true
  attribute :guest do |user|
    user.guest?
  end

  typelize :string, nullable: true
  attribute :name do |user|
    [ user.first_name, user.last_name ].join " "
  end

  typelize :boolean, nullable: true
  attribute :admin do |user|
    user.admin?
  end

  typelize :string, nullable: true
  attribute :password do |user| end

  typelize '"user"', nullable: false
  attribute :type do |user|
    "user"
  end
end
