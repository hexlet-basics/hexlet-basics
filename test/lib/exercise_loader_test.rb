require "test_helper"

class ExerciseLoaderTest < ActiveSupport::TestCase
  test "the truth" do
    version = language_versions(:javascript2)
    loader = ExerciseLoader.new
    loader.run(version)
  end
end
