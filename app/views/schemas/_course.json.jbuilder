# frozen_string_literal: true

json.set! :@context, 'https://schema.org'
json.set! :@type, 'Course'
json.name [language.current_version.name, language.learn_as.text].join(' ')
json.description language_version_info.description
json.accessMode 'textOnVisual'
# json.commentCount messages_count
json.url language_url(language.slug)
# json.image shrine_image_url(course.current_build.image, type: :fill, width: 600, height: 400, extension: :png, fallback: 'fallback/course.png')
json.isAccessibleForFree true
json.provider do
  json.partial! 'schemas/organization', format: [:json]
end

if language.members.finished.any?
  json.aggregateRating do
    json.set! :@type, 'AggregateRating'
    json.ratingValue "4.#{((language.members.finished.size / language.members.size) * 99).round}"
    json.reviewCount language.members.finished.size
  end
end
