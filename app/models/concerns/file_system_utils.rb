# frozen_string_literal: true

module FileSystemUtils
  def self.file_name_for_exercise(lesson, language)
    "#{lesson.id}.#{language.extension}"
  end

  def self.directory_for_code(user)
    chunked = user.id
                  .to_s
                  .rjust(6, '0')
                  .reverse
                  .split('')
                  .each_slice(3)
                  .to_a
                  .map(&:join)
    File.join(chunked)
  end
end
