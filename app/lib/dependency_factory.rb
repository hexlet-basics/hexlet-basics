# typed: strict

module DependencyFactory
  extend T::Sig

  sig { returns(Dependencies) }
  def self.build
    Dependencies.new(
      docker_exercise_client: Rails.env.test? ? DockerExerciseClientStub : DockerExerciseClient,
      google_one_tap: Rails.env.test? ? GoogleAuthStub : Google::Auth::IDTokens,
      amocrm: Amocrm::Client.new(),
      # Until a real SMS provider is chosen, the stub logs codes everywhere.
      # Swap to SmsSenderSmsc (or another adapter) for production when ready.
      sms_sender: SmsSenderStub
    )
  end
end
