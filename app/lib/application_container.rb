# frozen_string_literal: true

class ApplicationContainer
  extend Dry::Container::Mixin

  if Rails.env.test?
    register :docker_exercise_api, -> { DockerExerciseApiStub }
  else
    register :docker_exercise_api, -> { DockerExerciseApi }
  end
end
