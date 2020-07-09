class ExternalLinks
  LINKS = {
    hexlet_support_mail: 'support@hexlet.io',
    hexlet_facebook: 'https://www.facebook.com/codebasicsru',
    hexlet_telegram: 'https://t.me/hexlet_ru',
    hexlet_twitter: 'https://twitter.com/HexletHQ',
    hexlet_youtube: 'https://www.youtube.com/user/HexletUniversity',
    hexlet_slack: 'https://slack-ru.hexlet.io',
    hexlet_blog: 'https://ru.hexlet.io/blog',
    hexlet_recommended_books: 'https://ru.hexlet.io/pages/recommended-books',
    hexlet_matrix: 'https://docs.google.com/spreadsheets/d/e/2PACX-1vSr58Xg4DVX2WdmAkv3hE2ITJ3fPeNUmRFe0Ekro53U-ACFrhcUkV8PlUm4ddcn53Uh-5UIezZtZZgc/pubhtml',
    source_code: 'https://github.com/hexlet-basics',
  }

  def self.method_missing(method_name, *args, &block)
    LINKS[method_name]
  end
end
