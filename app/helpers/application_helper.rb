# frozen_string_literal: true

module ApplicationHelper
  include AuthConcern

  def let(value)
    yield value
  end

  def structured_data_tag(builder)
    content = builder.attributes!
    content['@context'] = 'https://schema.org'
    json = content.to_json
    tag.script(json.html_safe, type: 'application/ld+json')
  end

  def nav_menu_item(name, path = '#', options = {}, &block)
    assembled_options = options.merge(class: "nav-link #{options[:class]} #{active?(path)}".chomp)
    block_content = block ? capture(&block) : ''
    link = link_to name, path, assembled_options
    tag.li class: 'nav-item' do
      "#{link}#{block_content}".html_safe
    end
  end

  def active?(path, options = {})
    # raise options.inspect
    if options.key? :active_if
      'active' if options[:active_if]
    elsif current_page?(path)
      'active text-light'
    end
  end

  def markdown2html(text, options = {}, extensions = {})
    default_extensions = {
      autolink: true,
      tables: true,
      filter_html: false,
      safe_links_only: false,
      fenced_code_blocks: true,
      no_intra_emphasis: true
    }
    combined_extensions = default_extensions.merge(extensions)

    default_options = {
      no_images: false,
      escape_html: false,
      hard_wrap: true,
      prettify: true,
      with_toc_data: true,
      link_attributes: { target: '_blank' }
    }
    combined_options = default_options.merge options

    # TODO: use HtmlWithHoc
    # renderer = ::HtmlWithHoc.new(combined_options)
    # # renderer = Redcarpet::Render::HTML.new(combined_options)
    # markdown = Redcarpet::Markdown.new(renderer, combined_extensions)
    # markdown.render(text)

    html_toc = Redcarpet::Markdown.new(Redcarpet::Render::HTML_TOC)
    markdown = Redcarpet::Markdown.new(Redcarpet::Render::HTML.new(combined_options), combined_extensions)
    toc  = html_toc.render(text)
    html = markdown.render(text)
    "#{toc}#{html}"
  end

  def current_breadcrumb(lesson_version, lesson_version_info, langugage_lessons_count)
    "#{lesson_version_info} #{lesson_version.natural_order}/#{langugage_lessons_count}"
  end

  def get_lesson_source_code(lesson_version, lesson_version_info)
    repository_path = ExternalLinks.source_code_curl
    locale = lesson_version_info.locale
    path_to_description = File.join(repository_path, lesson_version.path_to_code, "description.#{locale}.yml")

    path_to_description.sub('modules', 'blob/main/modules')
  end

  def language_version_state_class(state)
    case state
    when 'created'
      'badge badge-secondary'
    when 'building'
      'badge badge-warning'
    when 'built'
      'badge badge-success'
    end
  end

  def get_language_background(slug)
    mapping = {
      'php' => 'x-bg-blue',
      'javascript' => 'x-bg-yellow',
      'java' => 'x-bg-orange',
      'python' => 'x-bg-orange',
      'html' => 'x-bg-orange',
      'css' => 'x-bg-azure',
      'racket' => 'x-bg-red',
      'ruby' => 'x-bg-red',
      'elixir' => 'x-bg-indigo',
      'go' => 'x-bg-cyan',
      'clang' => 'x-bg-azure',
      'clojure' => 'x-bg-green',
      'typescript' => 'x-bg-blue',
      'csharp' => 'x-bg-purple',
      'lua' => 'x-bg-blue',
      'prolog' => 'x-bg-orange',
      'haskell' => 'x-bg-gray',
      'cpp' => 'x-bg-blue',
      'bash' => 'x-bg-gray',
      'fortran' => 'x-bg-purple',
      'kotlin' => 'x-bg-purple',
      'swift' => 'x-bg-orange',
      'rust' => 'x-bg-blue',
      'perl' => 'x-bg-gray',
      'ocaml' => 'x-bg-blue'
    }
    mapping[slug]
  end

  def get_complete_icon_name(slug)
    mapping = {
      'php' => 'php.svg',
      'javascript' => 'frontend.svg',
      'java' => 'java.svg',
      'python' => 'python.svg',
      'html' => 'layout-designer.svg',
      'css' => 'layout-designer.svg',
      'clang' => 'clang.svg'
    }
    mapping.fetch(slug, 'hexlet_logo.png')
  end

  def get_language_devicon_name(slug)
    mapping = {
      'html' => 'html5',
      'css' => 'css3',
      'cpp' => 'cplusplus',
      'clang' => 'c',
      'racket' => 'devicon',
      'prolog' => 'devicon',
      'fortran' => 'devicon'
    }

    mapping.fetch(slug, slug)
  end

  def get_continue_study_path(slug)
    mapping = {
      'php' => ExternalLinks.hexlet_php_curl,
      'javascript' => ExternalLinks.hexlet_frontend_curl,
      'java' => ExternalLinks.hexlet_java_curl,
      'python' => ExternalLinks.hexlet_python_curl,
      'ruby' => ExternalLinks.hexlet_rails_curl,
      'html' => ExternalLinks.hexlet_layout_designer_curl,
      'css' => ExternalLinks.hexlet_layout_designer_curl
    }
    mapping.fetch(slug, ExternalLinks.hexlet_profession_curl)
  end

  def supported_browser?
    # NOTE https://github.com/browserslist/browserslist-useragent-ruby#3-add-helper
    # @browsers ||= JSON.parse(File.read('browsers.json'))
    # matcher = BrowserslistUseragent::Match.new(@browsers, request.user_agent)
    # matcher.browser? && matcher.version?(allow_higher: true)
    # TODO replace implementation
    true
  end

  def completed_languages
    @completed_languages ||= Language.with_progress(:completed).with_locale.ordered
  end

  def default_filter_form_options(options = {})
    { method: 'get', wrapper: :filter_input, html: { class: 'row row-cols-1 row-cols-sm-2 row-cols-lg-3 row-cols-xl-4 gx-2 align-items-end' }, url: url_for, defaults: { required: false, label: false } }.merge(options)
  end

  def language_menu_min_width_class(columns_count)
    case columns_count
    when 1
      'x-mw-160'
    when 2
      'x-mw-320'
    else
      raise "Invalid languages menu columns count: #{columns_count}"
    end
  end
end
