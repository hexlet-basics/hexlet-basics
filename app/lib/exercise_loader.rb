# frozen_string_literal: true

# rubocop:disable Metrics/ClassLength
class ExerciseLoader
  include Import['download_exercise_klass']

  def from_website(language_version)
    lang_name = language_version.language.slug
    language = language_version.language

    language_version.build!
    repo_dest = download_exercise_klass.run(lang_name)
    module_dest = "#{repo_dest}/modules"

    Language::Version.transaction do
      language.update!(version: language_version)

      update_language_version(repo_dest, language, language_version)

      modules_with_meta = get_modules(module_dest)
      language_modules = modules_with_meta.map { |data| find_or_create_module_with_data(language, data, language_version) }

      lessons = language_modules.flat_map { |language_module| get_lessons(module_dest, language_module, language) }
      lessons.each { |lesson| find_or_create_lesson_with_data_and_exercise(lesson, language_version) }

      language_version.update(result: 'Success')
      language_version.done!
    end
  rescue StandardError => e
    language_version.update(result: "Error class: #{e.class} message: #{e.message}")
    language_version.done!
  end

  def from_cli(lang_name)
    repo_dest = "tmp/hexletbasics/exercises-#{lang_name}"
    module_dest = "#{repo_dest}/modules"

    language = find_or_create_language(lang_name)
    language_version = create_language_version(repo_dest, language)
    language.update!(version: language_version)

    language_version.build!

    Language::Version.transaction do
      modules_with_meta = get_modules(module_dest)
      language_modules = modules_with_meta.map { |data| find_or_create_module_with_data(language, data, language_version) }

      lessons = language_modules.flat_map { |language_module| get_lessons(module_dest, language_module, language) }
      lessons.each { |lesson| find_or_create_lesson_with_data_and_exercise(lesson, language_version) }

      language_version.update(result: 'Success')
      language_version.done!
    end
  rescue StandardError => e
    language_version.update(result: "Error class: #{e.class} message: #{e.message}")
    language_version.done!
  end

  def get_modules(dest)
    files = Dir.glob("#{dest}/*")

    files
      .filter { |file| File.directory?(file) }
      .map do |directory|
        filename = File.basename(directory)
        order, slug = filename.split('-', 2)
        descriptions = get_descriptions(File.join(dest, filename))
        { order: order, slug: slug, descriptions: descriptions }
      end
  end

  def get_descriptions(path)
    files = Dir.glob("#{path}/description.*.yml")

    files.map do |file|
      filename = File.basename(file)
      _, locale, = filename.split('.')

      data = YAML.load_file(file)
      [locale, data]
    end
  end

  def get_lessons(dest, language_module, language)
    module_dir = "#{language_module.version.order}-#{language_module.slug}"
    module_path = File.join(dest, module_dir)
    wildcard_path = File.join(module_path, '*')
    files = Dir.glob(wildcard_path)

    files
      .filter { |file| File.directory?(file) }
      .map do |directory|
        filename = File.basename(directory)
        order, slug = filename.split('-', 2)

        descriptions = get_descriptions(directory)
        lesson_version = get_lesson_version(directory, language, language_module)

        {
          order: order,
          module: language_module,
          language: language,
          slug: slug,
          lesson_version: lesson_version,
          descriptions: descriptions
        }
      end
  end

  def get_lesson_version(directory, language, language_module)
    module_dir = "#{language_module.version.order}-#{language_module.slug}"
    test_file_path = File.join(directory, language.version.exercise_test_filename)
    test_code = File.read(test_file_path)
    original_code = File.read(File.join(directory, language.version.exercise_filename))
    prepared_code = prepare_code(original_code)
    path_to_code = File.join("/exercises-#{language.slug}/modules", module_dir, directory)

    {
      test_code: test_code,
      original_code: original_code,
      prepared_code: prepared_code,
      path_to_code: path_to_code
    }
  end

  def find_or_create_language(lang_name)
    Language.find_or_create_by!(slug: lang_name)
  end

  def update_language_version(repo_dest, language, language_version)
    spec_filepath = File.join(repo_dest, 'spec.yml')
    language_info = YAML.load_file(spec_filepath)['language']

    language_version.update!(
      name: language.slug,
      extension: language_info['extension'],
      docker_image: language_info['docker_image'],
      exercise_filename: language_info['exercise_filename'],
      exercise_test_filename: language_info['exercise_test_filename']
    )
  end

  def create_language_version(repo_dest, language)
    spec_filepath = File.join(repo_dest, 'spec.yml')
    language_info = YAML.load_file(spec_filepath)['language']

    Language::Version.create!(
      name: language.slug,
      extension: language_info['extension'],
      docker_image: language_info['docker_image'],
      exercise_filename: language_info['exercise_filename'],
      exercise_test_filename: language_info['exercise_test_filename'],
      language: language
    )
  end

  def find_or_create_module_with_data(language, data, language_version)
    order, slug, descriptions = data.values_at(:order, :slug, :descriptions)

    language_module = Language::Module.find_or_create_by!(slug: slug, language: language)

    version = Language::Module::Version.create!(
      order: order,
      language: language,
      language_version: language_version,
      module: language_module
    )

    language_module.update!(version: version)

    raise "Module: #{language.module} does not have data" if descriptions.empty?

    descriptions.each { |description| create_module_datum(language, language_version, language_module, description) }

    language_module
  end

  def create_module_datum(language, language_version, language_module, description_data)
    locale, info = description_data

    new_datum_attr = {
      module: language_module,
      language: language,
      language_version: language_version,
      locale: locale,
      version: language_module.version
    }.merge(info)

    Language::Module::Version::Datum.create!(new_datum_attr)
  end

  def find_or_create_lesson_with_data_and_exercise(data, language_version)
    language = data[:language]
    language_module = data[:module]
    slug = data[:slug]
    order = data[:order]
    descriptions = data[:descriptions]
    lesson_version = data[:lesson_version]

    lesson = Language::Lesson.find_or_create_by!(language: language, slug: slug, module: language_module)

    version = Language::Lesson::Version.create!(
      test_code: lesson_version[:test_code],
      order: order,
      original_code: lesson_version[:original_code],
      prepared_code: lesson_version[:prepared_code],
      path_to_code: lesson_version[:path_to_code],
      lesson: lesson,
      language_version: language_version,
      language: language,
      module_version: language_module.version
    )

    lesson.update!(version: version)
    raise "Lesson '#{language_module.slug}.#{lesson.slug}' does not have descriptions" if descriptions.empty?

    descriptions.each { |description| create_lesson_datum(language, language_version, version, lesson, description) }

    lesson
  end

  def create_lesson_datum(language, language_version, lesson_version, lesson, description_data)
    locale, info = description_data

    new_datum_attr = {
      locale: locale,
      lesson: lesson,
      language: language,
      language_version: language_version,
      version: lesson_version
    }.merge(info)

    Language::Lesson::Version::Datum.create!(new_datum_attr)
  end

  def prepare_code(code)
    reg = /(?<begin>^[^\n]*?BEGIN.*?$\s*)(?<content>.+?)(?<end>^[^\n]*?END.*?$)/msu

    result = code.gsub(reg, "\\1\n\\3")

    result != code ? result : ''
  end
end
# rubocop:enable Metrics/ClassLength
