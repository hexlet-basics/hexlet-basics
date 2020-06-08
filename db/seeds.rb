# frozen_string_literal: true
# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: 'Star Wars' }, { name: 'Lord of the Rings' }])
#   Character.create(name: 'Luke', movie: movies.first)
language = Language.create!(slug: 'JavaScript')
language_version = Language::Version.create!(
  docker_image: 'docker_image',
  exercise_filename: 'index.js',
  exercise_test_filename: 'test.js',
  extension: '.js',
  name: 'JavaScript',
  language_id: language.id,
)
language.update(current_version_id: language_version.id)


language_module = Language::Module.create!(
  slug: 'Basic',
  language_id: language.id,
)
language_module_version = Language::Module::Version.create!(
  language_version_id: language_version.id,
  module_id: language_module.id,
  order: 1,
)
language_module.update(current_version_id: language_module_version.id)


language_module_lesson = Language::Module::Lesson.create!(
  slug: 'variables',
  language_id: language.id,
  module_id: language_module.id,
)
language_module_lesson_version = Language::Module::Lesson::Version.create!(
  language_version_id: language_version.id,
  module_version_id: language_module_version.id,
  lesson_id: language_module_lesson.id,
  order: 1,
  original_code: 'original_code',
  prepared_code: 'prepared_code',
  test_code: 'test_code',
  path_to_code: 'path/to/code'
)
language_module_lesson.update(current_version_id: language_module_lesson_version.id)


language_module_description = Language::Module::Description.create!(
  name: 'Very loong Desctiprion name',
  description: 'Here is description on RU --- Faker.<>',
  locale: 'ru',
  module_id: language_module.id,
  language_id: language.id,
)
language_module_lesson_description = Language::Module::Lesson::Description.create!(
  instructions: 'Here is instructions',
  locale: 'ru',
  name: 'It is name, ',
  theory: 'Theory olo',
  tips: 'Tips, Tips, Tips',
  definitions: '- dedinition 1 -- definition 2',
  lesson_id: language_module_lesson.id,
  language_id: language.id,
)
