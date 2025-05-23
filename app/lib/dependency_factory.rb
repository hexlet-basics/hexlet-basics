# typed: strict

module DependencyFactory
  extend T::Sig

  sig { returns(Dependencies) }
  def self.build
    Dependencies.new(
      openai_api: OpenAI::Client.new,
      event_registry: EventRegistry.new,
      docker_exercise_client: Rails.env.test? ? DockerExerciseClientStub : DockerExerciseClient,
      google_one_tap: Rails.env.test? ? GoogleAuthStub : Google::Auth::IDTokens,
      n8n_client: Rails.env.test? ? N8nClientStub.new : N8nClient.new
    )
  end
end
