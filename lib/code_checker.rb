# frozen_string_literal: true

class CodeChecker
  def self.check(code, user, lesson_version, language_version)
    code_directory = Rails.configuration.hexlet_basics[:code_directory]
    full_directory_path = File.join(code_directory, FileSystemUtils.directory_for_code(user))
    FileUtils.mkdir_p(full_directory_path)

    full_exercise_file_path = File.join(full_directory_path, FileSystemUtils.file_name_for_exercise(lesson_version, language_version))
    File.open(full_exercise_file_path, 'w') { |f| f.write(code) }

    path_to_exersice_file = File.join(lesson_version.path_to_code, language_version.exercise_filename)
    volume = "-v #{full_exercise_file_path}:#{path_to_exersice_file}"

    docker_command_template = Rails.configuration.hexlet_basics[:docker_command_template]
    docker_command = format(docker_command_template, volume, language_version.docker_image, lesson_version.path_to_code)

    file_descriptor_path = File.join(full_directory_path, 'data.txt')
    file_descriptor = IO.sysopen(file_descriptor_path, 'w')
    io = IO.new(file_descriptor)

    system(docker_command, err: io, out: io)
    status = $CHILD_STATUS.exitstatus

    io_output_file = File.new(file_descriptor_path)
    output_content = io_output_file.read
    io_output_file.close

    { status: status, output: output_content }
  end
end
