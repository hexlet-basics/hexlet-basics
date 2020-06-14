# frozen_string_literal: true

class ExerciseLoaderWorker
  include Sidekiq::Worker

  def perform(upload_id)
    upload = Upload.find(upload_id)
    Exercises::Loader.run_with_upload(upload)
  end
end
