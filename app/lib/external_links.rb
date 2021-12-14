# frozen_string_literal: true

class ExternalLinks
  LINKS = {
    hexlet: 'https://hexlet.io/?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_support_mail: 'support@hexlet.io',
    hexlet_facebook: 'https://www.facebook.com/codebasicsru',
    hexlet_telegram: 'https://t.me/hexlet_ru',
    hexlet_twitter: 'https://twitter.com/HexletHQ',
    hexlet_youtube: 'https://www.youtube.com/user/HexletUniversity',
    hexlet_slack: 'https://slack-ru.hexlet.io',
    hexlet_blog: 'https://ru.hexlet.io/blog',
    hexlet_success_stories: 'https://ru.hexlet.io/blog/categories/success?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_recommended_books: 'https://ru.hexlet.io/pages/recommended-books',
    hexlet_matrix: 'https://docs.google.com/spreadsheets/d/e/2PACX-1vSr58Xg4DVX2WdmAkv3hE2ITJ3fPeNUmRFe0Ekro53U-ACFrhcUkV8PlUm4ddcn53Uh-5UIezZtZZgc/pubhtml',
    source_code: 'https://github.com/hexlet-basics',
    hexlet_profession: 'https://ru.hexlet.io/programs?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_layout_designer: 'https://ru.hexlet.io/programs/layout-designer?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_frontend: 'https://ru.hexlet.io/programs/frontend?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_java: 'https://ru.hexlet.io/programs/java?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_python: 'https://ru.hexlet.io/programs/python?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_php: 'https://ru.hexlet.io/programs/php?utm_campaign=general&utm_medium=referral&utm_source=code-basics',
    hexlet_rails: 'https://ru.hexlet.io/programs/rails?utm_campaign=general&utm_medium=referral&utm_source=code-basics'
  }.freeze

  def self.method_missing(method_name)
    return LINKS[method_name] if LINKS.key?(method_name)

    super
  end

  def self.respond_to_missing?(method_name)
    LINKS.key?(method_name)
  end
end
