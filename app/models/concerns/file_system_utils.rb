# frozen_string_literal: true

module FileSystemUtils
  def self.file_name_for_exercise(user, lesson_version, language)
    "#{user.id}.#{lesson_version.id}.#{language.extension}"
  end

  def self.directory_for_code(user)
    chunked = user.id
                  .to_s
                  .rjust(3, '0')
                  .reverse[0...3]
                  .chars
    File.join(chunked)
  end
end
