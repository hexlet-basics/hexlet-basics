# frozen_string_literal: true

namespace :exercises do
  desc 'Load exercies'
  task :load, [:lang] => :environment do |_task, args|
    ExerciseLoader.from_cli(args.lang)
  end
end
