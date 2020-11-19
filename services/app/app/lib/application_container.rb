# frozen_string_literal: true

class ApplicationContainer
  extend Dry::Container::Mixin

  if Rails.env.test?
    register :download_exercise_klass, -> { DownloadExerciseStub }
    register :open3_klass, -> { Open3Stub }
  else
    register :download_exercise_klass, -> { DownloadExercise }
    register :open3_klass, -> { Open3 }
  end
end
