# frozen_string_literal: true

module ApplicationHelper
  include AuthConcern

  # def let(value)
  #   yield value
  # end

  # def structured_data_tag(builder)
  #   content = builder.attributes!
  #   content['@context'] = 'https://schema.org'
  #   json = content.to_json
  #   tag.script(json.html_safe, type: 'application/ld+json')
  # end

  def get_lesson_source_code(lesson_version, lesson_version_info)
    repository_path = ExternalLinks.source_code_curl
    locale = lesson_version_info.locale
    path_to_description = File.join(repository_path, lesson_version.path_to_code, locale, 'README.md')

    path_to_description.sub('modules', 'blob/main/modules')
  end

  # def language_version_state_class(state)
  #   case state
  #   when 'created'
  #     'badge badge-secondary'
  #   when 'building'
  #     'badge badge-warning'
  #   when 'built'
  #     'badge badge-success'
  #   end
  # end

  # def get_complete_icon_name(slug)
  #   mapping = {
  #     'php' => 'php.svg',
  #     'javascript' => 'frontend.svg',
  #     'java' => 'java.svg',
  #     'python' => 'python.svg',
  #     'html' => 'layout-designer.svg',
  #     'css' => 'layout-designer.svg',
  #     'clang' => 'clang.svg',
  #     'layout-designer' => 'layout-designer.svg',
  #     'pre-course-java' => 'pre-course-java.svg',
  #     'pre-course-python' => 'pre-course-python.svg',
  #     'pre-course-javascript' => 'frontend.svg'
  #   }
  #   mapping.fetch(slug, 'hexlet_logo.png')
  # end

  # def get_continue_study_path(slug)
  #   mapping = {
  #     'php' => ExternalLinks.hexlet_php_curl,
  #     'javascript' => ExternalLinks.hexlet_frontend_curl,
  #     'java' => ExternalLinks.hexlet_java_curl,
  #     'python' => ExternalLinks.hexlet_python_curl,
  #     'ruby' => ExternalLinks.hexlet_rails_curl,
  #     'html' => ExternalLinks.hexlet_layout_designer_curl,
  #     'css' => ExternalLinks.hexlet_layout_designer_curl
  #   }
  #   mapping.fetch(slug, ExternalLinks.hexlet_profession_curl)
  # end

  # def completed_languages(locale = I18n.locale)
  #   Language.with_progress(:completed).with_locale(locale).ordered
  # end

  # def default_filter_form_options(options = {})
  #   { method: 'get', wrapper: :filter_input, html: { class: 'row row-cols-1 row-cols-sm-2 row-cols-lg-3 row-cols-xl-4 gx-2 align-items-end' }, url: url_for, defaults: { required: false, label: false } }.merge(options)
  # end

  # def language_menu_min_width_class(columns_count)
  #   case columns_count
  #   when 1
  #     'x-mw-160'
  #   when 2
  #     'x-mw-320'
  #   else
  #     raise "Invalid languages menu columns count: #{columns_count}"
  #   end
  # end
end
