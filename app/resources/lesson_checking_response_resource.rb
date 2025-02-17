class LessonCheckingResponseResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :passed, :output, :result, :status, :lesson_has_been_finished, :language_has_been_finished

  typelize passed: :boolean,
    output: :string,
    status: :number,
    lesson_has_been_finished: :boolean,
    language_has_been_finished: :boolean

    typelize result: [ enum: [ "error", "passed", "failed", "failed-infinity" ] ]
end
