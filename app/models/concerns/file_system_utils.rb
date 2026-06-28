# typed: strict
# frozen_string_literal: true

module FileSystemUtils
  extend T::Sig

  sig { params(user: T.untyped, lesson_version: T.untyped, language: T.untyped).returns(T.untyped) }
  def self.file_name_for_exercise(user, lesson_version, language)
    "#{user&.id}.#{lesson_version.id}.#{language.extension}"
  end

  sig { params(user: T.untyped).returns(T.untyped) }
  def self.directory_for_code(user)
    chunked = user&.id
                  .to_s
                  .rjust(3, "0")
                  .reverse[0...3]
                  .chars
    File.join(chunked)
  end
end
