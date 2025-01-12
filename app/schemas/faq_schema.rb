# frozen_string_literal: true

class FaqSchema
  class << self
    def to_builder
      Jbuilder.new do |json|
        json.set! :@type, "FAQPage"
        json.mainEntity faq.values.map do |item|
          json.set! :@type, "Question"
          json.name item.fetch(:question)
          json.acceptedAnswer do
            json.set! :@type, "Answer"
            json.text item.fetch(:answer)
          end
        end
      end
    end
  end
end
