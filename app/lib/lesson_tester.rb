# frozen_string_literal: true

class LessonTester
  include Import['docker_exercise_api']

  def run(lesson_version, language_version, code, user)
    code_directory = '/tmp/hexlet-basics/code'
    full_directory_path = File.join(code_directory, FileSystemUtils.directory_for_code(user))
    FileUtils.mkdir_p(full_directory_path)

    full_exercise_file_path = File.join(full_directory_path, FileSystemUtils.file_name_for_exercise(language_version, language_version))
    File.write(full_exercise_file_path, code || '')
    path_to_exersice_file = File.join(lesson_version.path_to_code, language_version.exercise_filename)

    volume = "-v #{full_exercise_file_path}:#{path_to_exersice_file}"
    command = "docker run --rm --net none #{volume} #{language_version.docker_image}:#{language_version.image_tag} timeout 4 make --silent -C #{lesson_version.path_to_code} test"

    output, process_status = docker_exercise_api.run_exercise(
      volume: volume,
      docker_image: language_version.docker_image,
      image_tag: language_version.image_tag,
      path_to_code: lesson_version.path_to_code
    )
    exitstatus = process_status.exitstatus

    result = case exitstatus
             when 0
               'passed'
             when 124
               'failed-infinity'
             else
               'failed'
             end
    passed = result == 'passed'

    { passed: passed, output: Base64.strict_encode64(output), result: result, status: exitstatus }
  end
end
