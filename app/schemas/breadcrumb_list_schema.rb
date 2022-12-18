# frozen_string_literal: true

class BreadcrumbListSchema
  class << self
    def to_builder(_items)
      Jbuilder.new do |json|
        json.set! '@type', 'BreadCrumbList'
        # json.itemListElement items.each_with_index.to_a do |(item, index)|
        #   json.set! '@type', 'ListItem'
        #   json.position index + 1
        #   json.item do
        #     json.set! '@id' item.fetch('url')
        #     json.name item.fetch('name')
        #   end
        # end
      end
    end
  end
end
