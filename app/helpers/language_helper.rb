# frozen_string_literal: true

module LanguageHelper
  def get_language_background(slug)
    mapping = {
      'php' => 'bg-blue',
      'javascript' => 'bg-yellow',
      'java' => 'bg-azure',
      'python' => 'bg-orange',
      'html' => 'bg-orange',
      'css' => 'bg-azure',
      'racket' => 'bg-red',
      'ruby' => 'bg-red',
      'elixir' => 'bg-indigo',
      'go' => 'bg-cyan'
    }
    mapping[slug]
  end
end
