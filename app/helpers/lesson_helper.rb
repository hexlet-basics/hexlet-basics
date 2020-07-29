# frozen_string_literal: true

module LessonHelper
  def markdown2html(text)
    extensions = {
      autolink: true,
      filter_html: false,
      safe_links_only: false,
      fenced_code_blocks: true
    }

    options = {
      escape_html: false,
      hard_wrap: true,
      prettify: true,
      link_attributes: { target: '_blank' }
    }

    renderer = Redcarpet::Render::HTML.new(options)
    markdown = Redcarpet::Markdown.new(renderer, extensions)
    markdown.render(text)
  end
end
