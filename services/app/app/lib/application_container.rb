# frozen_string_literal: true

class ApplicationContainer
  extend Dry::Container::Mixin

  if Rails.env.test?
    register :download_exercise_klass, -> { DockerExerciseApiStub }
    register :open3_klass, -> { Open3Stub }
  else
    register :download_exercise_klass, -> { DockerExerciseApi }
    register :open3_klass, -> { Open3 }
  end
end
