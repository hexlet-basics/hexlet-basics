# typed: strict

class BlogPostCreateResource < ApplicationResource
  typelize_from BlogPost

  attributes :state,
    :cover,
    :id,
    name: [ String, true ],
    slug: [ String, true ],
    description: [ String, true ],
    body: [ String, true ]

  typelize cover: "File | null"
  typelize id: [ :number, nullable: true ]
end
