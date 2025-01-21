class Language::MemberResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Member
  # root_key :user

  attributes :id, :user_id, :language_id, :state

  typelize state: [ enum: [ "started", "finished" ] ]

  typelize language_name: :string
  attribute :language_name do |member|
    member.language.name
  end
end
