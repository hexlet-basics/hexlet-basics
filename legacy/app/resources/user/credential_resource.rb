# typed: strict

class User::CredentialResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User::Credential

  attributes :id, :nickname, :created_at
end
