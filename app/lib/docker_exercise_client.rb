# frozen_string_literal: true

require "open3"

class DockerExerciseClient < DockerExerciseClientInterface
  def self.repo_dest(lang_name)
    "/tmp/hexletbasics/exercises-#{lang_name}"
  end

  def self.image_name(lang_name)
    "ghcr.io/hexlet-basics/exercises-#{lang_name}"
  end

  def self.download(lang_name)
    _stdout, stderr, status = Open3.capture3("docker pull #{image_name(lang_name)}")

    unless status.success?
      raise "Docker command failed: #{stderr.strip}"
    end

    system("rm -rf #{repo_dest(lang_name)}")

    # FIXME docker in docker volume
    system("docker run --name exercises-#{lang_name} -v #{repo_dest(lang_name)}:/out #{image_name(lang_name)}")
    system("docker cp exercises-#{lang_name}:/exercises-#{lang_name} /tmp/hexletbasics/")
    system("docker rm exercises-#{lang_name}")
  end

  def self.run_exercise(created_code_file_path:, exercise_file_path:, docker_image:, image_tag:, path_to_code:)
    volume = "-v #{created_code_file_path}:#{exercise_file_path}"
    command = "docker run --rm --memory=512m --network none #{volume} #{docker_image}:#{image_tag} timeout 6 make --silent -C #{path_to_code} test"
    Rails.logger.debug(command)

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [ output.join, status ]
  end

  def self.tag_image_version(lang_name)
    return unless Rails.env.production?

    tag_command = "docker tag #{image_name(lang_name)}:latest #{image_name(lang_name)}:release"
    raise "Docker tag error" unless system(tag_command)

    push_command = "docker push #{image_name(lang_name)}:release"
    raise "Docker push error" unless system(push_command)
  end
end
