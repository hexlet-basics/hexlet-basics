# typed: strict

class Dependencies < T::Struct
  const :docker_exercise_client, T.class_of(DockerExerciseClientInterface)
  const :google_one_tap, T.untyped
  const :amocrm, Amocrm::Client
  const :sms_sender, T.class_of(SmsSenderInterface)
end
