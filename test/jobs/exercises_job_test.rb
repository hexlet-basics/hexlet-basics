require 'test_helper'

class ExercisesJobTest < ActiveJob::TestCase
  test 'should create entities' do
    Rails.class_eval do
      def self.root
        File.join(Dir.pwd, 'test', 'fixtures', 'files')
      end
    end

    assert_difference -> { Language.count } => 1, -> { Language::Module.count } => 2, -> { Language::Module::Lesson.count } => 1 do
      ExercisesJob.perform_now('javascript')
    end
  end
end
