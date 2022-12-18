# frozen_string_literal: true

class ReviewSchema
  class << self
    # include Rails.application.routes.url_helpers

    def to_builder(review)
      info = review.language.current_version.infos.find_by!(locale: I18n.locale)

      Jbuilder.new do |json|
        json.set! '@type', 'Review'
        json.itemReviewed CourseSchema.to_builder(review.language, info)
        json.name review.body
        json.author do
          json.set! '@type', 'Person'
          json.name review.user.to_s
        end
      end
    end
  end
end
