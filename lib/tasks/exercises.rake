# frozen_string_literal: true

namespace :exercises do
  desc 'Load exercies'
  task :load, [:lang] => :environment do |_task, args|
    # remove from here
    language = Language.find_or_create_by!(slug: args.lang)
    language_version = language.versions.build
    language_version.save!

    ExerciseLoader.new.run(language_version)
  end
end
