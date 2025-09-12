# frozen_string_literal: true

class LessonTester
  def run(lesson_version, language_version, code, user)
    docker_exercise_client = DepsLocator.current.docker_exercise_client

    code_directory = "/tmp/hexlet-basics/code"
    full_directory_path = File.join(code_directory, FileSystemUtils.directory_for_code(user))
    FileUtils.mkdir_p(full_directory_path, verbose: false)

    created_code_file_path = File.join(full_directory_path, FileSystemUtils.file_name_for_exercise(user, lesson_version, language_version))
    File.write(created_code_file_path, code || "")
    Rails.logger.debug(created_code_file_path)
    exercise_file_path = File.join(lesson_version.path_to_code, language_version.exercise_filename)
    Rails.logger.debug(exercise_file_path)

    full_image_name = docker_exercise_client.ensure_image(
      image_name: language_version.docker_image,
      image_tag: language_version.image_tag,
      lang_name:  language_version.language.slug,
      lang_version: language_version.id
    )

    output, process_status = docker_exercise_client.run_exercise(
      created_code_file_path: created_code_file_path,
      exercise_file_path: exercise_file_path,
      full_image_name: full_image_name,
      path_to_code: lesson_version.path_to_code
    )
    exitstatus = process_status.exitstatus

    result = case exitstatus
    when 0
      "passed"
    when 124
      "failed-infinity"
    else
      "failed"
    end
    passed = result == "passed"

    # NOTE: escapeURIComponent нужен для удаления недопустимых символов в UTF-8, которые могут быть в выводе упражнения
    sanitized_output = result == "failed-infinity" ? "" : Base64.strict_encode64(CGI.escapeURIComponent(output))

    { passed: passed, output: sanitized_output, result: result, status: exitstatus }
  end
end
