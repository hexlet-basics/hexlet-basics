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

  def current_breadcrumb(lesson_version, lesson_version_info, langugage_lessons_count)
    "#{lesson_version_info} #{lesson_version.natural_order}/#{langugage_lessons_count}"
  end

  def get_lesson_source_code(lesson_version, lesson_version_info)
    repository_path = ExternalLinks.source_code
    locale = lesson_version_info.locale
    path_to_description = File.join(repository_path, lesson_version.path_to_code, "description.#{locale}.yml")

    path_to_description.gsub('modules', 'blob/master/modules')
  end
end
