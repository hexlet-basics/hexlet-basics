class ExercisesJob < ApplicationJob
  queue_as :default

  def perform(lang_name)
    exercises_path = File.join(Rails.root, 'tmp', 'hexletbasics')
    language_exercises_path = File.join(exercises_path, "exercises-#{lang_name}")
    language = upsert_language(language_exercises_path, lang_name)
    language_id = language.id

    modules_names = get_modules(language_exercises_path)

    modules_names.each do |module_name|
      modul = upsert_module(module_name, language_id)
      module_lessons_names = get_lessons(language_exercises_path, module_name)
      module_lessons_names.each do |module_lesson_name|
        upsert_lesson(module_lesson_name, modul.id, language_id)
      end
    end
  end

  def get_modules(dest)
    modules_path = File.join(dest, 'modules')
    Dir.open(modules_path) do |dir|
      children = dir.children
      modules_names = children.select { |child| File.directory?(File.join(modules_path, child)) }
    end
  end

  def get_lessons(dest, module_name)
    modules_path = File.join(dest, 'modules')
    lessons_path = File.join(modules_path, module_name)

    Dir.open(lessons_path) do |dir|
      children = dir.children
      lessons_names = children.select { |child| File.directory?(File.join(lessons_path, child)) }
    end
  end

  def upsert_lesson(lesson_name, module_id, language_id)
    Language::Module::Lesson.find_or_create_by(slug: lesson_name, language_module_id: module_id, language_id: language_id)
  end

  def upsert_module(module_name, language_id)
    Language::Module.find_or_create_by(slug: module_name, language_id: language_id)
  end

  def upsert_language(dest, lang_name)
    data_path = File.join(dest, 'spec.yml')
    data = File.new(data_path, 'r').read
    language_data = YAML.safe_load(data)['language']


    Language.find_or_create_by(
      name: lang_name,
      slug: lang_name,
      extension: language_data['extension']
    )
  end
end
