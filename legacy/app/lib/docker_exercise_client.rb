# typed: strict
# frozen_string_literal: true

require "open3"

class DockerExerciseClient < DockerExerciseClientInterface
  extend T::Sig

  sig { override.params(lang_name: String).returns(String) }
  def self.repo_dest(lang_name)
    "/tmp/hexletbasics/exercises-#{lang_name}"
  end

  sig { override.params(lang_name: String).returns(String) }
  def self.image_name(lang_name)
    "ghcr.io/hexlet-basics/exercises-#{lang_name}"
  end

  sig { override.params(lang_name: String).void }
  def self.download(lang_name)
    run_command!("docker pull #{image_name(lang_name)}")

    system("rm -rf #{repo_dest(lang_name)}")
    system("mkdir -p #{repo_dest(lang_name)}")

    # FIXME docker in docker volume
    run_command!("docker run --name exercises-#{lang_name} -v #{repo_dest(lang_name)}:/out #{image_name(lang_name)}")
    run_command!("docker cp exercises-#{lang_name}:/exercises-#{lang_name} /tmp/hexletbasics/")
    system("docker rm exercises-#{lang_name}")
  end

  sig { override.params(created_code_file_path: String, exercise_file_path: String, full_image_name: String, path_to_code: String).returns([ String, T.untyped ]) }
  def self.run_exercise(created_code_file_path:, exercise_file_path:, full_image_name:, path_to_code:)
    volume = "-v #{created_code_file_path}:#{exercise_file_path}"
    command = "docker run --rm --memory=512m --memory-swap=-1 --network none #{volume} #{full_image_name} timeout 6 make --silent -C #{path_to_code} test"
    Rails.logger.debug(command)

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [ output.join, status ]
  end

  sig { override.params(lang_name: String).void }
  def self.tag_image_version(lang_name)
    return unless Rails.env.production?

    tag_command = "docker tag #{image_name(lang_name)}:latest #{image_name(lang_name)}:release"
    raise "Docker tag error" unless system(tag_command)

    push_command = "docker push #{image_name(lang_name)}:release"
    raise "Docker push error" unless system(push_command)
  end

  sig { override.params(image_name: String, image_tag: String, lang_name: String, lang_version: Integer).returns(String) }
  def self.ensure_image(image_name:, image_tag:, lang_name:, lang_version:)
    full_image_name = "#{image_name(lang_name)}:v#{lang_version}"

    unless system("docker image inspect #{full_image_name} > /dev/null 2>&1")
      system("docker pull #{image_name}:#{image_tag}")
      system("docker tag #{image_name}:#{image_tag} #{full_image_name}")
    end

    full_image_name
  end

  sig { params(cmd: String).void }
  def self.run_command!(cmd)
    Rails.logger.debug(cmd)
    _stdout, stderr, status = Open3.capture3(cmd)

    unless status.success?
      raise "Docker command failed: #{cmd} (#{stderr.strip})"
    end
  end
end
