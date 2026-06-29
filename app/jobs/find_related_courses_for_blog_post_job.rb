# typed: strict

class FindRelatedCoursesForBlogPostJob < ApplicationJob
  extend T::Sig

  sig { params(blog_post_id: Integer).void }
  def perform(blog_post_id)
    blog_post = BlogPost.find(blog_post_id)
    landing_pages = Language::LandingPage.published.where(main: true).with_locale(blog_post.locale).includes(:language)

    instructions = <<~PROMPT
      Ты — ассистент, который помогает подобрать курсы.
      У тебя есть текст статьи блога и список курсов.
      Выбери пять курсов подходящих под тему статьи в порядке приоритета. Первый наиболее близок, последний - наименее.
      Верни результат в виде JSON-массива идентификаторов курсов (по полю `id`) отсортированный по похожести.
      Первыми должны идти наиболее близкие курсы.
    PROMPT

    languages_data = landing_pages.map do |lp|
      { id: lp.language.id, name: lp.header }
    end

    content = blog_post.body
    response = RubyLLM.chat
      .with_instructions(instructions)
      .ask([
        "Текст статьи: #{content.truncate(2000)}",
        "Список курсов: #{languages_data.to_json}"
      ].join("\n\n"))

    raw_output = response.content.to_s.gsub(/\A```(?:json)?|```\z/, "").strip
    course_ids = JSON.parse(raw_output)

    Rails.logger.info "COURSES #{course_ids}"

    return if course_ids.empty?

    blog_post.related_language_items.delete_all
    # languages = Language.where(id: course_ids)

    course_ids.each_with_index do |id, i|
      item = blog_post.related_language_items.build
      item.language_id = id
      item.order = i + 1
      item.save!
    end
  end
end
