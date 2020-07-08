# frozen_string_literal: true

class ExerciseLoaderJob < ApplicationJob
  def perform(language_version_id)
    version = Language::Version.find(language_version_id)

    ExerciseLoader.new.run(version)
  end
end
