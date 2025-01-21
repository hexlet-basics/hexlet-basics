class Language::Lesson::MemberResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Member

  attributes :id, :user_id, :state

  typelize state: [ enum: [ "started", "finished" ] ]
end
