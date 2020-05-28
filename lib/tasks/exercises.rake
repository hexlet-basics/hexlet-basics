# frozen_string_literal: true

# frozen_string_litera: true

namespace :exercies do
  desc 'Loader exercies'
  task :load, [:lang] => :environment do |_task, args|
    Exercises::Loader.new(args.lang).run
  end
end
