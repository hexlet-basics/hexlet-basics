# typed: strict

module FakerExt
  extend T::Sig

  sig { returns(String) }
  def self.markdown
    [
      Faker::Markdown.headers,
      Faker::Markdown.emphasis,
      Faker::Markdown.sandwich(sentences: 6, repeat: 3),
      Faker::Markdown.headers,
      Faker::Markdown.block_code,
      Faker::Markdown.headers,
      Faker::Markdown.table,
      Faker::Markdown.headers,
      Faker::Markdown.ordered_list
    ].join("\n\n")
  end
end
