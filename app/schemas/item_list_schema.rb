# frozen_string_literal: true

class ItemListSchema
  class << self
    def to_builder(item_builders)
      Jbuilder.new do |json|
        json.set! '@type', 'ItemList'
        json.itemListElement item_builders.each_with_index.to_a do |(item, index)|
          json.set! '@type', 'ListItem'
          json.position index + 1
          json.item item

          attributes = item.attributes!

          json.url attributes.fetch('url')
          json.name attributes['name']
        end
      end
    end
  end
end
