# frozen_string_literal: true

class DockerExerciseApiStub
  def self.repo_dest(_)
    'test/fixtures/files/exercises'
  end

  def self.download(_); end

  def self.tag_image_version(_lang_version, _tag); end

  # rubocop:disable Lint/UnusedMethodArgument
  def self.run_exercise(created_code_file_path:, exercise_file_path:, docker_image:, image_tag:, path_to_code:)
    command = "ruby #{Rails.root.join('test/fixtures/files/exercise/test.rb')} #{created_code_file_path}"

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [output.join, status]
  end

  def self.remove_image(lang_name, tag); end
  # rubocop:enable Lint/UnusedMethodArgument
end
