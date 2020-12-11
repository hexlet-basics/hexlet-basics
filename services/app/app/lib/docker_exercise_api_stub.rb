# frozen_string_literal: true

class DockerExerciseApiStub
  def self.download(_)
    'test/fixture_files/exercises'
  end

  def self.tag_image_version(_lang_version, _tag); end
end
