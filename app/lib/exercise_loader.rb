# frozen_string_literal: true

# rubocop:disable Metrics/ClassLength
class ExerciseLoader
  include Import['download_exercise_klass']

  def run(language_version)
    lang_name = language_version.language.slug
    language = language_version.language

    language_version.build!

    repo_dest = download_exercise_klass.run(lang_name)
    module_dest = "#{repo_dest}/modules"

    update_language_version(repo_dest, language, language_version)

    modules_with_meta = get_modules(module_dest)
    language_modules = modules_with_meta.map { |data| find_or_create_module_with_info(language, data, language_version) }

    lessons = language_modules.flat_map { |language_module| get_lessons(module_dest, language_module, language_version) }
    lessons.each { |lesson| find_or_create_lesson_with_info_and_exercise(lesson) }

    language_version.update(result: 'Success')
    language_version.done!
    language.update!(current_version: language_version)
  rescue StandardError => e
    language_version.update(result: "Error class: #{e.class} message: #{e.message}")
    language_version.done!
    raise
  end

  private

  def get_modules(dest)
    files = Dir.glob("#{dest}/*")

    files
      .filter { |file| File.directory?(file) }
      .map do |directory|
        filename = File.basename(directory)
        order, slug = filename.split('-', 2)
        infos = get_infos(File.join(dest, filename))
        { order: order, slug: slug, infos: infos }
      end
  end

  def get_infos(path)
    files = Dir.glob("#{path}/description.*.yml")

    files.map do |file|
      filename = File.basename(file)
      _, locale, = filename.split('.')

      data = YAML.load_file(file)
      [locale, data]
    end
  end

  def get_lessons(dest, language_module, language_version)
    module_dir = "#{language_module.current_version.order}-#{language_module.slug}"
    module_path = File.join(dest, module_dir)
    wildcard_path = File.join(module_path, '*')
    files = Dir.glob(wildcard_path)

    files
      .filter { |file| File.directory?(file) }
      .map do |directory|
        filename = File.basename(directory)
        order, slug = filename.split('-', 2)

        infos = get_infos(directory)
        lesson_version = get_lesson_version(directory, language_version, language_module)

        {
          order: order,
          module: language_module,
          language: language_version.language,
          language_version: language_version,
          slug: slug,
          lesson_version: lesson_version,
          infos: infos
        }
      end
  end

  def get_lesson_version(directory, language_version, language_module)
    module_dir = "#{language_module.current_version.order}-#{language_module.slug}"
    test_file_path = File.join(directory, language_version.exercise_test_filename)
    test_code = File.read(test_file_path)
    original_code = File.read(File.join(directory, language_version.exercise_filename))
    prepared_code = prepare_code(original_code)
    path_to_code = File.join("/exercises-#{language_version.language.slug}/modules", File.basename(module_dir), File.basename(directory))

    {
      test_code: test_code,
      original_code: original_code,
      prepared_code: prepared_code,
      path_to_code: path_to_code
    }
  end

  def update_language_version(repo_dest, language, language_version)
    spec_filepath = File.join(repo_dest, 'spec.yml')
    language_info = YAML.load_file(spec_filepath).fetch('language')

    language_version.assign_attributes(
      name: language.slug,
      extension: language_info['extension'],
      docker_image: language_info['docker_image'],
      exercise_filename: language_info['exercise_filename'],
      exercise_test_filename: language_info['exercise_test_filename']
    )

    language.save!

    language
  end

  def find_or_create_module_with_info(language, data, language_version)
    order, slug, infos = data.values_at(:order, :slug, :infos)

    language_module = Language::Module.find_or_create_by!(slug: slug, language: language)

    version = Language::Module::Version.new(
      order: order,
      language: language,
      language_version: language_version,
      module: language_module
    )

    version.save!
    language_module.update!(current_version: version)

    raise "Module: #{language.module} does not have info" if infos.empty?

    infos.each { |info| create_module_info(language, language_version, language_module, info) }

    language_module
  end

  def create_module_info(language, language_version, language_module, info_data)
    locale, data = info_data

    new_datum_attr = {
      language: language,
      language_version: language_version,
      locale: locale,
      version: language_module.current_version
    }.merge(data)

    info = Language::Module::Version::Info.new(new_datum_attr)
    info.save!

    info
  end

  def find_or_create_lesson_with_info_and_exercise(data)
    language = data[:language]
    language_module = data[:module]
    slug = data[:slug]
    order = data[:order]
    infos = data[:infos]
    lesson_version = data[:lesson_version]
    language_version = data[:language_version]

    lesson = Language::Lesson.find_or_create_by!(language: language, slug: slug, module: language_module)

    version = Language::Lesson::Version.new(
      test_code: lesson_version[:test_code],
      order: order,
      original_code: lesson_version[:original_code],
      prepared_code: lesson_version[:prepared_code],
      path_to_code: lesson_version[:path_to_code],
      lesson: lesson,
      language_version: language_version,
      language: language,
      module_version: language_module.current_version
    )

    version.save!
    lesson.update!(current_version: version)
    raise "Lesson '#{language_module.slug}.#{lesson.slug}' does not have info" if infos.empty?

    infos.each { |info| create_lesson_info(language, language_version, version, info) }

    lesson
  end

  def create_lesson_info(language, language_version, lesson_version, info_data)
    locale, data = info_data

    new_datum_attr = {
      locale: locale,
      language: language,
      language_version: language_version,
      version: lesson_version
    }.merge(data)

    info = Language::Lesson::Version::Info.new(new_datum_attr)
    info.save!

    info
  end

  def prepare_code(code)
    reg = /(?<begin>^[^\n]*?BEGIN.*?$\s*)(?<content>.+?)(?<end>^[^\n]*?END.*?$)/msu

    result = code.gsub(reg, "\\1\n\\3")

    result != code ? result : ''
  end
end
# rubocop:enable Metrics/ClassLength
