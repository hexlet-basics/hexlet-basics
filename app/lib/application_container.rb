# frozen_string_literal: true

class ApplicationContainer
  extend Dry::Container::Mixin

  register :openai_api, -> { OpenAI::Client.new }
  register :event_registry, -> { EventRegistry.new }

  if Rails.env.test?
    register :docker_exercise_api, -> { DockerExerciseApiStub }
    register :google_one_tap, -> { GoogleAuthStub }
    register :n8n_client, -> { N8nClientStub.new }
  else
    register :docker_exercise_api, -> { DockerExerciseApi }
    register :google_one_tap, -> { Google::Auth::IDTokens }
    register :n8n_client, -> { N8nClient.new }
  end
end
