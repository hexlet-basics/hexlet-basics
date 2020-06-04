# frozen_string_literal: true

namespace :exercises do
  desc 'Load exercies'
  task :load, [:lang] => :environment do |_task, args|
    Exercises::Loader.run(args.lang)
  end
end
