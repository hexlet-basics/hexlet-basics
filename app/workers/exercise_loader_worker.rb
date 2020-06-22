# frozen_string_literal: true

class ExerciseLoaderWorker
  include Sidekiq::Worker

  def perform(upload_id)
    upload = Language::Upload.find(upload_id)

    ExerciseLoader.new.from_website(upload)
  end
end
