# typed: strict

class ExerciseLoaderJob < ApplicationJob
  extend T::Sig

  sig { params(language_version_id: Integer).void }
  def perform(language_version_id)
    version = Language::Version.find(language_version_id)

    ExerciseLoader.new.run(version)
  end
end
