# frozen_string_literal: true

namespace :data_utils do
  task :create_missing_language_members, [:limit] => :environment do |_task, args|
    limit = args.limit || 1000

    puts 'Processing'

    lang_lesson_members = Language::Lesson::Member.joins('LEFT JOIN language_members ON language_members.user_id = language_lesson_members.user_id AND language_members.language_id = language_lesson_members.language_id')
                                                  .where(language_members: { id: nil })
                                                  .select('language_lesson_members.user_id, language_lesson_members.language_id')
                                                  .group('language_lesson_members.user_id, language_lesson_members.language_id')
                                                  .limit(limit)
                                                  .pluck('language_lesson_members.user_id, language_lesson_members.language_id')

    puts "Total processing: #{lang_lesson_members.size}"

    lang_lesson_members.each do |user_id, lang_id|
      puts "Process: #{user_id}, #{lang_id}"

      user = User.find user_id
      lang = Language.find lang_id

      language_member = user.language_members.build(language: lang)

      if language_member.may_finish?
        language_member.finish!
      else
        language_member.save!
      end
    end

    puts 'Finish'
  end

  task :populate_language_members, [:limit] => :environment do |_task, args|
    limit = args.limit || 1000

    puts 'Processing'

    lang_lesson_members = Language::Lesson::Member.left_joins(:language_member)
                                                  .where(language_members: { id: nil })
                                                  .limit(limit)

    lang_lesson_members.find_each do |lesson_member|
      puts "Process: #{lesson_member.id}"

      user = lesson_member.user
      lang = lesson_member.language

      language_member = user.language_members.find_by!(language: lang)

      lesson_member.language_member = language_member
      lesson_member.save!
    end

    puts 'Finish'
  end
end
