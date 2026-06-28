# typed: strict

class MarkdownImageProcessor
  extend T::Sig

  sig { params(markdown: String, block: T.proc.params(url: String).returns(String)).returns(String) }
  def self.process(markdown, &block)
    markdown.gsub(/!\[([^\]]*)\]\(([^)]+)\)/) do |match|
      alt_text = Regexp.last_match(1)
      image_url = Regexp.last_match(2)

      unless image_url
        next
      end

      # Delegate processing logic to a block
      new_path = !image_url.start_with?("http") ? yield(image_url) : image_url

      "![#{alt_text}](#{new_path})"
    end
  end
end
