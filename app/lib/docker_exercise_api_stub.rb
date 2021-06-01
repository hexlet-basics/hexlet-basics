# frozen_string_literal: true

class DockerExerciseApiStub
  def self.repo_dest(_)
    'test/fixture_files/exercises'
  end

  def self.download(_); end

  def self.tag_image_version(_lang_version, _tag); end

  # rubocop:disable Lint/UnusedMethodArgument
  def self.run_exercise(created_code_file_path:, path_to_exersice_file:, docker_image:, image_tag:, path_to_code:)
    command = "ruby #{Rails.root.join('test/fixtures/files/exercise/test.rb')} #{created_code_file_path}"

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [output.join, status]
  end
  # rubocop:enable Lint/UnusedMethodArgument
end
