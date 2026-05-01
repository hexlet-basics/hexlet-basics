class BlogPostCreateResource < ApplicationResource
  typelize_from BlogPost

  attributes :state,
    :cover,
    id: [ Integer, true ],
    name: [ String, true ],
    slug: [ String, true ],
    description: [ String, true ],
    body: [ String, true ]

  typelize cover: "File | null"
end
