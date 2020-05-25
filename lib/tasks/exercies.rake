# frozen_string_literal: true

# frozen_string_litera: true

namespace :exercies do
  desc 'Pull exercies'
  task :pull, [:lang] => :environment do |_task, args|
    ExercisesLoadJob.perform_now(args.lang)
  end
end
