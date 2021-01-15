# frozen_string_literal: true

namespace :exercises do
  desc 'Load exercies'
  task :load, [:lang] => :environment do |_task, args|
    language_version = LanguageVersionManager.new.find_or_create_language_with_version(args.lang)

    ExerciseLoader.new.run(language_version)
  end
end
