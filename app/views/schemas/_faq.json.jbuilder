# frozen_string_literal: true

json.set! :@context, 'https://schema.org'
json.set! :@type, 'FAQPage'
json.mainEntity faq.values.map do |item|
  json.set! :@type, 'Question'
  json.name item.fetch(:question)
  json.acceptedAnswer do
    json.set! :@type, 'Answer'
    json.text item.fetch(:answer)
  end
end
