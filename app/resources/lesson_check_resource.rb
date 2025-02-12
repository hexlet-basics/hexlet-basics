class LessonCheckResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :attributes, :passing_of_entities

  typelize attributes: "Record<'output', string | 'passed', boolean | 'result', string | 'status', number>",
    passing_of_entities: "Record<'is_lesson_become_finished', boolean | 'is_lesson_become_finished', boolean>"
end
