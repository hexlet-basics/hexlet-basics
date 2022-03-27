# frozen_string_literal: true

class DockerExerciseApi
  def self.repo_dest(lang_name)
    "/tmp/#{image_name(lang_name)}"
  end

  def self.image_name(lang_name)
    "hexletbasics/exercises-#{lang_name}"
  end

  def self.download(lang_name)
    system("docker pull #{image_name(lang_name)}")
    system("rm -rf #{repo_dest(lang_name)}")

    # FIXME docker in docker volume
    system("docker run --name exercises-#{lang_name} -v #{repo_dest(lang_name)}:/out #{image_name(lang_name)}")
    system("docker cp exercises-#{lang_name}:/exercises-#{lang_name} /tmp/hexletbasics/")
    system("docker rm exercises-#{lang_name}")
  end

  def self.run_exercise(created_code_file_path:, exercise_file_path:, docker_image:, image_tag:, path_to_code:)
    volume = "-v #{created_code_file_path}:#{exercise_file_path}"
    command = "docker run --rm --net none #{volume} #{docker_image}:#{image_tag} timeout 5 make --silent -C #{path_to_code} test"
    Rails.logger.debug(command)

    output = []
    status = BashRunner.start(command) { |line| output << line }

    [output.join, status]
  end

  def self.tag_image_version(lang_name, tag)
    return unless Rails.env.production?

    tag_command = "docker tag #{image_name(lang_name)}:latest #{image_name(lang_name)}:#{tag}"
    BashRunner.start(tag_command)

    push_command = "docker push #{image_name(lang_name)}:#{tag}"
    ok = BashRunner.start(push_command)

    raise "Docker tag error: #{ok}" unless ok
  end

  def self.remove_image(lang_name, tag)
    remove_command = "docker rmi #{image_name(lang_name)}:#{tag}"
    _ok = BashRunner.start(remove_command)
  end
end
