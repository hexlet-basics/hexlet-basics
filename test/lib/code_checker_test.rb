# frozen_string_literal: true

require 'test_helper'

class CodeCheckerTest < ActionDispatch::IntegrationTest
  test 'should check a code' do
    Rails.class_eval do
      def self.configuration
        config = {
          code_directory: Dir.tmpdir,
          docker_command_template: 'echo'
        }
        OpenStruct.new(hexlet_basics: config)
      end
    end

    lesson_version = language_lesson_versions(:one)
    language_version = lesson_version.language_version

    user = users(:one)
    code_data = 'code'

    dest_path = File.join(Dir.tmpdir, FileSystemUtils.directory_for_code(user))
    file_path = File.join(dest_path, FileSystemUtils.file_name_for_exercise(lesson_version, language_version))

    result = CodeChecker.check(code_data, user, lesson_version, language_version)

    File.open(file_path) do |file|
      assert { file.read == code_data }
    end

    assert { (result.keys == %i[status output]) }
    assert { result[:status].zero? }
  end
end
