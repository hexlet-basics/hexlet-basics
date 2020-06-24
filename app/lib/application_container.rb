# frozen_string_literal: true

class ApplicationContainer
  extend Dry::Container::Mixin

  if Rails.env.test?
    register :download_exercise_klass, -> { :success }
  else
    register :download_exercise_klass, -> { DownloadExercise }
  end
end
