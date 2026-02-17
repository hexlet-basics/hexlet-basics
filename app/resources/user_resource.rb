class UserResource < ApplicationResource
  typelize_from User

  attributes :id,
    :first_name,
    :email,
    :created_at,
    :last_name,
    :admin,
    :assistant_messages_count

  typelize :boolean, nullable: true
  attribute :guest do
    it.guest?
  end

  typelize :string, nullable: true
  attribute :name do
    [ it.first_name, it.last_name ].join " "
  end

  typelize :boolean, nullable: true
  attribute :admin do
    it.admin?
  end

  typelize :string, nullable: true
  attribute :password do end

  # For Wootric (GTM)
  typelize :number, nullable: true
  attribute :created_at_as_timestamp do
    it.created_at.to_i
  end

  typelize '"user"', nullable: false
  attribute :type do
    "user"
  end
end
