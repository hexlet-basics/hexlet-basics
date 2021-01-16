# frozen_string_literal: true

class ExerciseLoaderJob < ApplicationJob
  sidekiq_options retry: 0

  def perform(language_version_id)
    version = Language::Version.find(language_version_id)

    ExerciseLoader.new.run(version)
  end
end
