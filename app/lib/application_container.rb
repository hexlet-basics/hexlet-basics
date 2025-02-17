# frozen_string_literal: true

class ApplicationContainer
  extend Dry::Container::Mixin

  if Rails.env.test?
    register :docker_exercise_api, -> { DockerExerciseApiStub }
    register :google_one_tap, -> { GoogleAuthStub }
  else
    register :docker_exercise_api, -> { DockerExerciseApi }
    register :google_one_tap, -> { Google::Auth::IDTokens }
  end
end
