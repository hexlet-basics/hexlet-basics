class Dependencies < T::Struct
  const :openai_api, OpenAI::Client
  const :event_registry, EventRegistry
  const :docker_exercise_client, T.class_of(DockerExerciseClientInterface)
  const :google_one_tap, T.untyped
  const :amocrm, Amocrm::Client
end
