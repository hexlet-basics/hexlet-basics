class FindRelatedCoursesForBlogPostJob < ApplicationJob
  def perform(blog_post_id)
    blog_post = BlogPost.find(blog_post_id)
    landing_pages = Language::LandingPage.where(main: true).with_locale(blog_post.locale)

    openai_api = DepsLocator.current.openai_api

    instructions = <<~PROMPT
      Ты — ассистент, который помогает подобрать курсы.
      У тебя есть текст статьи блога и список курсов.
      Найди 5 курсов, которые наиболее соответствуют содержанию статьи.
      Верни результат в виде JSON-массива идентификаторов курсов (по полю `id`) отсортированный по похожести.
      Первыми должны идти наиболее близкие курсы.
    PROMPT

    languages_data = landing_pages.map do |lp|
      { id: lp.language&.id, name: lp.header, description: lp.description }
    end

    content = blog_post.body || ""
    chat_completion = openai_api.chat(
      parameters: {
        model: :"gpt-4.1",
        messages: [
          { role: "system", content: instructions },
          { role: "user", content: "Текст статьи: #{content.truncate(2000)}" },
          { role: "user", content: "Список курсов: #{languages_data.to_json}" }
        ]
      }
    )

    raw_output = chat_completion.dig("choices", 0, "message", "content")
    course_ids = JSON.parse(raw_output)
    # raise [ languages_data, course_ids ].inspect

    logger.info "[Job #{job_id}] Courses #{course_ids}"
    return if course_ids.empty?

    blog_post.related_language_items.clear
    languages = Language.where(id: course_ids)

    languages.each do |language|
      item = blog_post.related_language_items.build language: language
      item.save!
    end
  end
end
