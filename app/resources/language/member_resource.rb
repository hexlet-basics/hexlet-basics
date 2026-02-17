class Language::MemberResource < ApplicationResource
  typelize_from Language::Member
  # TODO: make it nullable
  # one :next_lesson_info, resource: Language::LessonResource
  # has_one :language, resource: LanguageResource
  # root_key :user

  attributes :id, :user_id, :language_id, :state

  # typelize state: [ enum: [ "started", "finished" ] ]

  typelize :number
  attribute :progress do
    lessons_count = it.language.current_lessons.size
    finished_lessons_count = it.lesson_members.finished.size
    if lessons_count != 0
      result = ((finished_lessons_count.to_f / lessons_count) * 100).ceil
      # NOTE: Такая ситуация возможна, когда в новой версии курса уроков меньше
      # TODO: В идеале надо учитывать пройденные уроки только через current_lessons
      result > 100 ? 100 : result
    else
      0
    end
  end

  typelize :string, nullable: true
  attribute :next_lesson_name do
    it.next_lesson_info&.name
  end
end
