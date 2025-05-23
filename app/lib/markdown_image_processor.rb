class MarkdownImageProcessor
  def self.process(markdown)
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
