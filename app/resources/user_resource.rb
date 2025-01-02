class UserResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User

  attributes :id, :first_name, :email

  typelize :boolean, nullable: true
  attribute :guest do |user|
    user.guest?
  end

  typelize :string, nullable: true
  attribute :password do |user| end

  typelize '"user"', nullable: false
  attribute :type do |user|
    "user"
  end
end
