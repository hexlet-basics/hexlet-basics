# frozen_string_literal: true

# Set the host name for URL creation
SitemapGenerator::Sitemap.create_index = true
SitemapGenerator::Sitemap.max_sitemap_links = 45_000

if Rails.env.production?
  SitemapGenerator::Sitemap.adapter = SitemapGenerator::AwsSdkAdapter.new(
    configus.sitemap.bucket.name,
    **configus.sitemap.bucket.credentials
  )

  SitemapGenerator::Sitemap.compress = true
else
  SitemapGenerator::Sitemap.compress = false
end

module SitemapGeneratorHelper
  def build_alternates(options)
    existed_in_locales = I18n.available_locales.filter do |locale|
      options[:current] == locale || options[:check_exists].call(locale)
    end

    alternates = existed_in_locales.map do |locale|
      { href: options[:url].call(AppHost.locale_for_url(locale)), lang: locale }
    end

    if alternates.size == 1
      alternates << { href: options[:url].call(AppHost.locale_for_url(options[:current])), lang: :'x-default' }
    elsif alternates.size > 1
      alternates << { href: options[:url].call(:ru), lang: :'x-default' }
    end

    alternates
  end
end

SitemapGenerator::Interpreter.include SitemapGeneratorHelper

I18n.available_locales.each do |current_locale|
  I18n.with_locale(current_locale) do
    SitemapGenerator::Sitemap.sitemaps_path = "sitemaps/#{current_locale}/"
    SitemapGenerator::Sitemap.default_host = root_url

    SitemapGenerator::Sitemap.create do
      group(filename: :languages) do
        scope = Language.with_progress(:completed).joins(current_version: :infos)
        scope.merge(Language::Version::Info.with_locale).find_each do |language|
          alternates = build_alternates(current: current_locale,
                                        check_exists: ->(locale) { scope.merge(Language::Version::Info.with_locale(locale)).exists?(slug: language.slug) },
                                        url: ->(locale) { language_url(id: language.slug, locale: locale) })
          add language_path(language.slug, locale: AppHost.locale_for_url(current_locale)), changefreq: :weekly, lastmod: nil, alternates: alternates, priority: 0.9

          language.lessons.find_each do |lesson|
            alternates = build_alternates(current: current_locale,
                                          check_exists: ->(locale) { scope.merge(Language::Version::Info.with_locale(locale)).find_by(slug: language.slug)&.lessons&.exists?(slug: lesson.slug) },
                                          url: ->(locale) { language_lesson_url(language.slug, lesson.slug, locale: locale) })
            add language_lesson_path(language.slug, lesson.slug, locale: AppHost.locale_for_url(current_locale)), changefreq: :weekly, lastmod: nil, priority: 0.8, alternates: alternates
          end
        end
      end
    end
  end
end
