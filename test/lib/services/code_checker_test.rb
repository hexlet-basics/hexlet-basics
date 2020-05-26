# frozen_string_literal: true

require 'test_helper'
require 'services/code_checker'

class Service::CodeCheckerTest < ActionDispatch::IntegrationTest
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

    lesson = Language.first.lessons.first
    user = users(:one)
    code_data = 'code'

    dest_path = File.join(Dir.tmpdir, user.directory_for_code)
    file_path = File.join(dest_path, lesson.file_name_for_exercise)

    result = Service::CodeChecker.check(code_data, user, lesson)

    File.open(file_path) do |file|
      assert_equal(file.read, code_data)
    end

    assert_equal(result.keys, %i[status output])
    assert_equal(result[:status], 0)
  end
end
