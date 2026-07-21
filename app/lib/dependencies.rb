# typed: strict

class Dependencies < T::Struct
  const :docker_exercise_client, T.class_of(DockerExerciseClientInterface)
  # Либо GoogleAuthStub (класс), либо Google::Auth::IDTokens (gem-модуль) —
  # оба вызываются как `.verify_oidc(...)`. Общий надтип обоих значений — модуль.
  const :google_one_tap, T::Module[T.anything]
  const :amocrm, Amocrm::Client
  const :sms_sender, T.class_of(SmsSenderInterface)
end
