# frozen_string_literal: true

class DockerExerciseApi
  def self.repo_dest(lang_name)
    "/var/tmp/hexletbasics/exercises-#{lang_name}"
  end

  def self.download(lang_name)
    system("docker pull hexletbasics/exercises-#{lang_name}")
    system("rm -rf /var/tmp/hexletbasics/exercises-#{lang_name}")

    # FIXME docker in docker volume
    system("docker run --name exercises-#{lang_name} -v #{repo_dest(lang_name)}:/out hexletbasics/exercises-#{lang_name}")
    system("docker cp exercises-#{lang_name}:/exercises-#{lang_name} /var/tmp/hexletbasics/")
    system("docker rm exercises-#{lang_name}")
  end

  def self.run_exercise(created_code_file_path:, exercise_file_path:, docker_image:, image_tag:, path_to_code:)
    volume = "-v #{created_code_file_path}:#{exercise_file_path}"
    command = "docker run --rm --net none #{volume} #{docker_image}:#{image_tag} timeout 4 make --silent -C #{path_to_code} test"

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [output.join, status]
  end

  def self.tag_image_version(lang_name, tag)
    return unless Rails.env.production?

    tag_command = "docker tag hexletbasics/exercises-#{lang_name}:latest hexletbasics/exercises-#{lang_name}:#{tag}"
    BashRunner.start(tag_command)

    push_command = "docker push hexletbasics/exercises-#{lang_name}:#{tag}"
    _ok = BashRunner.start(push_command)

    # FIXME better error handling
    # raise "Docker tag error: #{ok}" unless ok
  end
end
