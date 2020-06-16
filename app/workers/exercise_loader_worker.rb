# frozen_string_literal: true

class ExerciseLoaderWorker
  include Sidekiq::Worker

  def perform(language_id, upload_id)
    language = Language.find(language_id)
    upload = Language::Upload.find(upload_id)

    ExerciseLoader.from_website(language, upload)
  end
end
