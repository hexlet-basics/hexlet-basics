# frozen_string_literal: true

class ExerciseLoaderJob < ApplicationJob
  def perform(version_id)
    version = Language::Version.find(version_id)

    ExerciseLoader.new.from_website(version)
  end
end
