# frozen_string_literal: true

module LessonHelper
  def markdown_to_html(text)
    Kramdown::Document.new(text).to_html.html_safe
  end
end
