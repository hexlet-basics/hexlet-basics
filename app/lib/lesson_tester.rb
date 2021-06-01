# frozen_string_literal: true

class LessonTester
  include Import['docker_exercise_api']

  def run(lesson_version, language_version, code, user)
    code_directory = '/tmp/hexlet-basics/code'
    full_directory_path = File.join(code_directory, FileSystemUtils.directory_for_code(user))
    FileUtils.mkdir_p(full_directory_path)

    created_code_file_path = File.join(full_directory_path, FileSystemUtils.file_name_for_exercise(language_version, language_version))
    File.write(created_code_file_path, code || '')
    exercise_file_path = File.join(lesson_version.path_to_code, language_version.exercise_filename)

    output, process_status = docker_exercise_api.run_exercise(
      created_code_file_path: created_code_file_path,
      exercise_file_path: exercise_file_path,
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
