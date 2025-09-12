# frozen_string_literal: true

class DockerExerciseClientStub < DockerExerciseClientInterface
  def self.repo_dest(_)
    "test/fixtures/files/exercises"
  end

  def self.download(_); end

  def self.tag_image_version(_lang_version); end

  # rubocop:disable Lint/UnusedMethodArgument
  def self.run_exercise(created_code_file_path:, exercise_file_path:, full_image_name:, path_to_code:)
    command = "ruby #{Rails.root.join('test/fixtures/files/exercise/test.rb')} #{created_code_file_path}"

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [ output.join, status ]
  end

  def self.image_name(lang_name); end

  def self.ensure_image(image_name:, image_tag:, lang_name:, lang_version:); end
end
