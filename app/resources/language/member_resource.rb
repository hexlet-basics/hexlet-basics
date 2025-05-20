class Language::MemberResource
  include Alba::Resource
  include Typelizer::DSL

  has_one :next_lesson_info, key: :next_lesson, resource: Language::LessonForListsResource

  typelize_from Language::Member
  # root_key :user

  attributes :id, :user_id, :language_id, :state

  # typelize state: [ enum: [ "started", "finished" ] ]

  typelize :number
  attribute :progress do |obj|
    lessons_count = obj.language.current_lessons.size
    finished_lessons_count = obj.lesson_members.finished.size
    if lessons_count != 0
      result = ((finished_lessons_count.to_f / lessons_count) * 100).ceil
      # NOTE: Такая ситуация возможна, когда в новой версии курса уроков меньше
      # TODO: В идеале надо учитывать пройденные уроки только через current_lessons
      result > 100 ? 100 : result
    else
      0
    end
  end
end
