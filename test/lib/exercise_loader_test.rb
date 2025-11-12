require "test_helper"

class ExerciseLoaderTest < ActiveSupport::TestCase
  def test_the_truth
    version = language_versions(:javascript2)
    loader = ExerciseLoader.new
    loader.run(version)
  end
end
