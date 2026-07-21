# typed: strict
# frozen_string_literal: true

class DockerExerciseClientStub < DockerExerciseClientInterface
  extend T::Sig

  sig { override.params(lang_name: String).returns(String) }
  def self.repo_dest(lang_name)
    "test/fixtures/files/exercises"
  end

  sig { override.params(lang_name: String).void }
  def self.download(lang_name); end

  sig { override.params(lang_name: String).void }
  def self.tag_image_version(lang_name); end

  # rubocop:disable Lint/UnusedMethodArgument
  sig { override.params(created_code_file_path: String, exercise_file_path: String, full_image_name: String, path_to_code: String).returns([ String, T.untyped ]) }
  def self.run_exercise(created_code_file_path:, exercise_file_path:, full_image_name:, path_to_code:)
    command = "ruby #{Rails.root.join('test/fixtures/files/exercise/test.rb')} #{created_code_file_path}"

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [ output.join, status ]
  end

  sig { override.params(lang_name: String).returns(String) }
  def self.image_name(lang_name)
    ""
  end

  sig { override.params(image_name: String, image_tag: String, lang_name: String, lang_version: Integer).returns(String) }
  def self.ensure_image(image_name:, image_tag:, lang_name:, lang_version:)
    ""
  end
end
