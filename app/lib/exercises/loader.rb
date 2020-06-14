# frozen_string_literal: true

# rubocop:disable Metrics/ClassLength
class Exercises::Loader
  class << self
    def run_with_upload(upload)
      lang_name = upload.language_name
      system("docker pull hexletbasics/exercises-#{lang_name}")
      system("rm -rf tmp/hexletbasics/exercises-#{lang_name}")
      system("docker run --rm -v #{Dir.pwd}/tmp/hexletbasics/exercises-#{lang_name}:/out hexletbasics/exercises-#{lang_name} bash -c 'cp -r /exercises-#{lang_name}/* /out'")

      repo_dest = "tmp/hexletbasics/exercises-#{lang_name}"
      module_dest = "#{repo_dest}/modules"

      Upload.transaction do
        language = find_or_create_language_with_version(repo_dest, lang_name, upload)

        modules_with_meta = get_modules(module_dest)
        language_modules = modules_with_meta.map { |data| find_or_create_module_with_descriptions(language, data, upload) }

        lessons = language_modules.flat_map { |language_module| get_lessons(module_dest, language_module, language) }
        lessons.each { |lesson| find_or_create_lesson_with_descriptions_and_exercise(lesson, upload) }
      end
    end

    def run(lang_name)
      repo_dest = "tmp/hexletbasics/exercises-#{lang_name}"
      module_dest = "#{repo_dest}/modules"

      Upload.transaction do
        upload = Upload.create!(language_name: lang_name)
        language = find_or_create_language_with_version(repo_dest, lang_name, upload)

        modules_with_meta = get_modules(module_dest)
        language_modules = modules_with_meta.map { |data| find_or_create_module_with_descriptions(language, data, upload) }

        lessons = language_modules.flat_map { |language_module| get_lessons(module_dest, language_module, language) }
        lessons.each { |lesson| find_or_create_lesson_with_descriptions_and_exercise(lesson, upload) }
      end
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
      module_dir = "#{language_module.current_version.order}-#{language_module.slug}"
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
      module_dir = "#{language_module.current_version.order}-#{language_module.slug}"
      test_file_path = File.join(directory, language.current_version.exercise_test_filename)
      test_code = File.read(test_file_path)
      original_code = File.read(File.join(directory, language.current_version.exercise_filename))
      prepared_code = prepare_code(original_code)
      path_to_code = File.join("/exercises-#{language.slug}/modules", module_dir, directory)

      {
        test_code: test_code,
        original_code: original_code,
        prepared_code: prepared_code,
        path_to_code: path_to_code
      }
    end

    def find_or_create_language_with_version(repo_dest, lang_name, upload)
      spec_filepath = File.join(repo_dest, 'spec.yml')

      language_info = YAML.load_file(spec_filepath)['language']

      language = Language.find_or_create_by!(slug: lang_name)

      version = Language::Version.create!(
        name: lang_name,
        extension: language_info['extension'],
        docker_image: language_info['docker_image'],
        exercise_filename: language_info['exercise_filename'],
        exercise_test_filename: language_info['exercise_test_filename'],
        language: language,
        upload: upload
      )

      language.update!(current_version: version)

      language
    end

    def find_or_create_module_with_descriptions(language, data, upload)
      order, slug, descriptions = data.values_at(:order, :slug, :descriptions)

      language_module = Language::Module.find_or_create_by!(slug: slug, language: language)

      version = Language::Module::Version.create!(
        order: order,
        language_version: language.current_version,
        module: language_module,
        upload: upload
      )

      language_module.update!(current_version: version)

      raise "Module: #{language.module} does not have descriptions" if descriptions.empty?

      descriptions.each { |description| upsert_module_description(language_module, description) }

      language_module
    end

    def upsert_module_description(language_module, description_data)
      locale, data = description_data

      description = Language::Module::Description.find_or_initialize_by(
        module: language_module,
        language: language_module.language,
        locale: locale
      )

      new_data = { language: language_module.language }.merge(data)
      description.update!(new_data)

      description
    end

    def find_or_create_lesson_with_descriptions_and_exercise(data, upload)
      language = data[:language]
      language_module = data[:module]
      slug = data[:slug]
      order = data[:order]
      descriptions = data[:descriptions]
      lesson_version = data[:lesson_version]

      lesson = Language::Module::Lesson.find_or_create_by!(language: language, module: language_module, slug: slug)

      version = Language::Module::Lesson::Version.create!(
        test_code: lesson_version[:test_code],
        order: order,
        original_code: lesson_version[:original_code],
        prepared_code: lesson_version[:prepared_code],
        path_to_code: lesson_version[:path_to_code],
        lesson: lesson,
        language_version: language.current_version,
        module_version: language_module.current_version,
        upload: upload
      )

      lesson.update!(current_version: version)

      raise "Lesson '#{language_module.slug}.#{lesson.slug}' does not have descriptions" if descriptions.empty?

      descriptions.each { |description| upsert_lesson_description(lesson, description) }

      lesson
    end

    def upsert_lesson_description(lesson, description_data)
      locale, data = description_data

      description = Language::Module::Lesson::Description.find_or_initialize_by(
        lesson: lesson,
        language: lesson.language,
        locale: locale
      )

      new_data = { lesson: lesson }.merge(data)

      description.update!(new_data)
      description
    end

    def prepare_code(code)
      reg = /(?<begin>^[^\n]*?BEGIN.*?$\s*)(?<content>.+?)(?<end>^[^\n]*?END.*?$)/msu

      result = code.gsub(reg, "\\1\n\\3")

      result != code ? result : ''
  end
end
end
# rubocop:enable Metrics/ClassLength
