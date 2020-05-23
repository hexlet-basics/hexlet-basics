namespace :exercies do
  desc "Pull exercies"
  task :pull, [:lang] => :environment do |task, args|
    lang_name = args.lang
    Exercies.new(lang_name).run
  end
end

require 'yaml'

class Exercies
  attr_reader :lang_name, :logger

  def initialize(lang_name, logger = Logger.new(STDOUT))
    @lang_name = lang_name
    @logger = logger
  end

  def run
    repo_dest = "tmp/hexletbasics/exercises-#{lang_name}"
    module_dest = "#{repo_dest}/modules"

    language = upsert_language(repo_dest)
    modules_with_meta = get_modules(module_dest)

    modules = modules_with_meta { |data| upsert_modules_with_description(language, data) }
  end

  def get_modules(dest)
    files = Dir.glob("#{dest}/*")
    dirs = files.filter { |file| File.directory?(file) }

    dirs.map do |dir|
      filename = File.basename(dir)
      # TODO: add logger
      order, slug = filename.split('-')
      descriptions = get_descriptions(File.join(dest, filename))
      { order: order, slug: slug, descriptions: descriptions }
    end
  end

  def get_descriptions(path)
    files = Dir.glob("#{path}/description.*.yml")

    files.map do |file|
      filename = File.basename(file)
      _, locale, _ = filename.split('.')
      logger.debug file

      data = YAML.load_file(file)
      { locale: locale, data: data }
    end
  end

  def upsert_language(repo_dest)
    spec_filepath = File.join(repo_dest, "spec.yml")

    language_info = YAML.load_file(spec_filepath)['language']

    language = Language.find_or_create_by(slug: lang_name)

    language.update(
      name: lang_name,
      extension: language_info['extension']
    )

    language
  end

  def upsert_module_with_description(language, data)
    order = data['order']
    slug = data['slug']
    descriptions = data['descriptions']

    language_module = Language::Module.find_or_create_by(slug: slug, language: language)

    language_module.update(
      order: order
    )

    raise "Module: #{language.module} does not have descriptions" if descriptions.empty?
  end

  def upsert_module_description(language_module, description)

  end
end
