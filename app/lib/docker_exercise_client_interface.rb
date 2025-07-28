# typed: strict

class DockerExerciseClientInterface
  extend T::Sig
  extend T::Helpers

  abstract!

  sig { abstract.params(lang_name: String).returns(String) }
  def self.repo_dest(lang_name); end

  sig { abstract.params(lang_name: String).returns(String) }
  def self.image_name(lang_name); end

  sig { abstract.params(lang_name: String).void }
  def self.download(lang_name); end

  sig do
    abstract.params(
      created_code_file_path: String,
      exercise_file_path: String,
      docker_image: String,
      image_tag: String,
      path_to_code: String
    ).returns([ String, T.untyped ])
  end
  def self.run_exercise(created_code_file_path:, exercise_file_path:, docker_image:, image_tag:, path_to_code:); end

  sig { abstract.params(lang_name: String).void }
  def self.tag_image_version(lang_name); end
end
