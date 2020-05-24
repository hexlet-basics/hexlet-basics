# frozen_string_literal: true

class ExercisesJob < ApplicationJob
  queue_as :default

  def perform(lang_name)
    exercises_path = Rails.root.join('tmp/hexletbasics')
    language_exercises_path = File.join(exercises_path, "exercises-#{lang_name}")
    language = upsert_language(language_exercises_path, lang_name)

    modules_dest = File.join(language_exercises_path, 'modules')
    modules_data = get_modules(modules_dest)

    modules = modules_data.map do |module_data|
      upsert_module(language, module_data)
    end

    modules
      .flat_map { |language_module| get_lessons(modules_dest, language_module, language) }
      .each { |lesson_data| upsert_lesson(lesson_data) }
  end

  def get_modules(dest)
    Dir.open(dest) do |dir|
      children = dir.children
      children
        .select { |child| File.directory?(File.join(dest, child)) }
        .map do |module_name|
          order, slug = module_name.split('-', 2)
          { order: order, slug: slug }
        end
    end
  end

  def get_lessons(dest, language_module, language)
    module_directory = Language::Module.get_directory(language_module)
    module_path = File.join(dest, module_directory)

    Dir.open(module_path) do |dir|
      children = dir.children
      children
        .select { |child| File.directory?(File.join(module_path, child)) }
        .map do |lesson_name|
          order, slug = lesson_name.split('-', 2)
          path_to_code = File.join("/exercises-#{language.slug}/modules", module_directory, lesson_name)

          { order: order, slug: slug, path_to_code: path_to_code, language: language, language_module: language_module }
        end
    end
  end

  def upsert_lesson(data)
    Language::Module::Lesson.find_or_create_by(slug: data[:slug], language_module_id: data[:language_module].id, language_id: data[:language].id, order: data[:order], path_to_code: data[:path_to_code])
  end

  def upsert_module(language, data)
    Language::Module.find_or_create_by(slug: data[:slug], language_id: language.id, order: data[:order])
  end

  def upsert_language(dest, lang_name)
    data_path = File.join(dest, 'spec.yml')
    data = File.new(data_path, 'r').read
    language_data = YAML.safe_load(data)['language']

    Language.find_or_create_by(
      name: lang_name,
      slug: lang_name,
      extension: language_data['extension'],
      docker_image: language_data['docker_image'],
      exercise_filename: language_data['exercise_filename'],
      exercise_test_filename: language_data['exercise_test_filename']
    )
  end
end
