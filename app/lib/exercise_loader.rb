# frozen_string_literal: true

class ExerciseLoader
  include Import['docker_exercise_api']

  def run(language_version)
    return unless language_version.may_build?

    language_version.build!

    language = language_version.language
    lang_name = language.slug

    docker_exercise_api.download(lang_name)

    update_language(language_version)

    language_modules_data = create_modules(language_version)

    create_lessons(language_version, language_modules_data)

    # FIXME we should stop building image if docker answer code is not 200
    docker_exercise_api.tag_image_version(lang_name, language_version.image_tag)

    # TODO: rename to building_error_descriptoin and use only for error messages
    language_version.result = 'Success'
    ActiveRecord::Base.transaction do
      language_version.mark_as_built!
      language_version.language.update!(current_version: language_version)
    end
  rescue StandardError => e
    language_version.update(result: "Error class: #{e.class} message: #{e.message}")
    language_version.mark_as_failed!
  end

  private

  def create_modules(language_version)
    module_dest = "#{docker_exercise_api.repo_dest(language_version.language.slug)}/modules"

    modules_with_meta = get_modules(module_dest).sort_by { |language_module| language_module[:order] }
    modules_with_meta.map { |module_meta| create_module_hierachy(language_version, module_meta) }
  end

  def create_lessons(language_version, language_modules_data)
    module_dest = "#{docker_exercise_api.repo_dest(language_version.language.slug)}/modules"

    lessons = language_modules_data.flat_map do |module_data|
      unordered_lessons = get_lessons(module_dest, module_data[:module_version], language_version)
      unordered_lessons.sort_by { |lesson| lesson[:order] }
    end
    lessons.each_with_index { |lesson, index| create_lesson_hierarchy(lesson, index) }
  end

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

  def get_lessons(dest, module_version, language_version)
    language_module = module_version.module
    module_dir = "#{module_version.order}-#{language_module.slug}"
    module_path = File.join(dest, module_dir)
    wildcard_path = File.join(module_path, '*')
    files = Dir.glob(wildcard_path)

    files
      .filter { |file| File.directory?(file) }
      .map do |directory|
        filename = File.basename(directory)
        order, slug = filename.split('-', 2)

        infos = get_infos(directory)
        lesson_version = get_lesson_version(directory, language_version, module_version)

        {
          order: order,
          module: language_module,
          language: language_version.language,
          language_version: language_version,
          slug: slug,
          lesson_version: lesson_version,
          infos: infos,
          module_version: module_version
        }
      end
  end

  def get_lesson_version(directory, language_version, module_version)
    language_module = module_version.module
    module_dir = "#{module_version.order}-#{language_module.slug}"

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

  def update_language(language_version)
    repo_dest = docker_exercise_api.repo_dest(language_version.language.slug)
    spec_filepath = File.join(repo_dest, 'spec.yml')
    language_spec = YAML.load_file(spec_filepath).fetch('language')

    language_version.assign_attributes(
      name: language_version.language.slug,
      extension: language_spec['extension'],
      docker_image: language_spec['docker_image'],
      exercise_filename: language_spec['exercise_filename'],
      exercise_test_filename: language_spec['exercise_test_filename']
    )

    language_version.save!

    infos = get_infos(repo_dest)
    infos.each do |_locale, info_spec|
      language_version_info = language_version.infos.build(language: language_version.language)
      language_version_info.description = info_spec.fetch('description')
      language_version_info.save!
    end
  end

  def create_module_hierachy(language_version, data)
    order, slug, infos = data.values_at(:order, :slug, :infos)
    language = language_version.language

    language_module = Language::Module.find_or_create_by!(slug: slug, language: language)

    version = Language::Module::Version.new(
      order: order,
      language: language,
      language_version: language_version,
      module: language_module
    )

    version.save!

    raise "Module: #{language_module.slug} does not have info" if infos.empty?

    module_infos = infos.map { |info| create_module_info(language_version, version, info) }

    { language_module: language_module, module_version: version, module_infos: module_infos }
  end

  def create_module_info(language_version, module_version, info_data)
    locale, data = info_data
    language = language_version.language

    new_datum_attr = {
      language: language,
      language_version: language_version,
      locale: locale,
      version: module_version
    }.merge(data)

    info = Language::Module::Version::Info.new(new_datum_attr)
    info.save!

    info
  end

  def create_lesson_hierarchy(data, index)
    language = data[:language]
    language_module = data[:module]
    module_version = data[:module_version]
    slug = data[:slug]
    order = data[:order]
    infos = data[:infos]
    lesson_version = data[:lesson_version]
    language_version = data[:language_version]

    lesson = Language::Lesson.find_or_initialize_by(language: language, slug: slug)
    lesson.module = language_module

    raise "Lesson Validation error #{lesson.errors.inspect}" unless lesson.valid?

    lesson.save!

    version = Language::Lesson::Version.new(
      test_code: lesson_version[:test_code],
      order: order,
      original_code: lesson_version[:original_code],
      prepared_code: lesson_version[:prepared_code],
      path_to_code: lesson_version[:path_to_code],
      lesson: lesson,
      language_version: language_version,
      language: language,
      module_version: module_version,
      natural_order: index + 1
    )

    version.save!

    raise "Lesson '#{language_module.slug}.#{lesson.slug}' does not have info" if infos.empty?

    lesson_infos = infos.map { |info| create_lesson_info(language_version, version, info) }

    { lesson: lesson, lesson_version: version, lesson_infos: lesson_infos }
  end

  def create_lesson_info(language_version, lesson_version, info_data)
    locale, data = info_data

    new_datum_attr = {
      locale: locale,
      language: language_version.language,
      language_version: language_version,
      version: lesson_version
    }.merge(data)

    info = Language::Lesson::Version::Info.new(new_datum_attr)
    info.save!

    info
  end

  def prepare_code(code)
    reg = /(?<begin>^[^\n]*?BEGIN.*?$\s*)(?<content>.+?)(?<end>^[^\n]*?END.*?$)/msu

    result = code.gsub(reg, "\\k<begin>\n\\k<end>")
    template_code = result.gsub('BEGIN', 'BEGIN (write your solution here)')

    result == code ? '' : template_code
  end
end
