namespace :lesson_version_info do
  desc 'Populate Language::Lesson::Version::Info model with given locale'
  task :populate, [:locale] => :environment do |task, args|
    locale = args[:locale]

    if locale.nil? || locale.empty?
      puts 'Please provide a locale. Example: rake lesson_version_info:populate["es"]'
      return
    end

    Language::Lesson::Version::Info.where(locale: locale).destroy_all if locale != "en"

    en_records = Language::Lesson::Version::Info.where(locale: 'en')

    puts 'No records found with locale "en"' and return if en_records.empty?

    en_records.each do |en_record|
      new_record = Language::Lesson::Version::Info.new(
        language_id: en_record.language_id,
        language_version_id: en_record.language_version_id,
        version_id: en_record.version_id,
        locale: locale,
        name: en_record.name,
        description: en_record.description,
        theory: en_record.theory,
        tips: en_record.tips,
        definitions: en_record.definitions,
        instructions: en_record.instructions,
      )

      if new_record.save(validate: false)
        puts "Record successfully created with locale \"#{locale}\""
      else
        puts "Error creating record with locale \"#{locale}\":"
        puts new_record.errors.full_messages
      end
    end
  end
end