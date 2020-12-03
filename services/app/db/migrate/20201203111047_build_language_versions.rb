class BuildLanguageVersions < ActiveRecord::Migration[6.0]
  def change
    languages = Language.all
    languages.each do |language|
      @version = language.versions.build
      @version.save!
      ExerciseLoaderJob.perform_later(@version.id)
    end
  end
end
