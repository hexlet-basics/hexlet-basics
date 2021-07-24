# frozen_string_literal: true

# Add your own tasks in files placed in lib/tasks ending in .rake,
# for example lib/tasks/capistrano.rake, and they will automatically be available to Rake.

require_relative 'config/application'

Rails.application.load_tasks

# https://github.com/rails/webpacker/issues/405
Rake::Task['yarn:install'].clear
namespace :yarn do
  task :install do
    # Redefine as empty
  end
end
