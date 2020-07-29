# frozen_string_literal: true

module LanguageHelper
  def get_language_background(slug)
    mapping = {
      'php' => 'x-bg-blue',
      'javascript' => 'x-bg-yellow',
      'java' => 'x-bg-azure',
      'python' => 'x-bg-orange',
      'html' => 'x-bg-orange',
      'css' => 'x-bg-azure',
      'racket' => 'x-bg-red',
      'ruby' => 'x-bg-red',
      'elixir' => 'x-bg-indigo',
      'go' => 'x-bg-cyan'
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
      'css' => 'layout-designer.svg'
    }
    mapping.fetch(slug, 'hexlet_logo.png')
  end

  def get_continue_study_path(slug)
    mapping = {
      'php' => ExternalLinks.hexlet_php,
      'javascript' => ExternalLinks.hexlet_frontend,
      'java' => ExternalLinks.hexlet_java,
      'python' => ExternalLinks.hexlet_python,
      'html' => ExternalLinks.hexlet_layout_designer,
      'css' => ExternalLinks.hexlet_layout_designer
    }
    mapping.fetch(slug, ExternalLinks.hexlet_profession)
  end
end
