# frozen_string_literal: true

class ItemListSchema
  class << self
    def to_builder(items)
      Jbuilder.new do |json|
        json.set! '@type', 'ItemList'
        json.array! items.each_with_index.to_a do |(item, index)|
          json.set! '@type', 'ListItem'
          json.position index + 1
          json.item item
        end
      end
    end
  end
end
