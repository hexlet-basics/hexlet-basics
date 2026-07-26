# typed: true

require "test_helper"

class ExerciseLoaderTest < ActiveSupport::TestCase
  def test_run_builds_version
    version = language_versions(:javascript2)
    language = version.language

    assert { version.created? }
    assert { version.may_build? }

    loader = ExerciseLoader.new
    loader.run(version)

    language.reload

    assert { version.result == "Success" }
    assert { version.built? }
    assert { language.current_version == version }
  end

  def test_run_records_skip_when_version_is_not_buildable
    version = language_versions(:javascript1)
    version.build!

    loader = ExerciseLoader.new
    result = loader.run(version)

    version.reload

    assert { result == false }
    assert { version.result&.start_with?("Skipped") }
  end
end
