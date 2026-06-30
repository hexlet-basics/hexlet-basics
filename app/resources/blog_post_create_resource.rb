# typed: strict

class BlogPostCreateResource < ApplicationResource
  typelize_from BlogPost

  attributes :state,
    :cover,
    :id,
    name: [ String, true ],
    slug: [ String, true ],
    description: [ String, true ]

  typelize cover: "File | null"
  typelize id: [ :number, nullable: true ]

  typelize :string
  attribute :rich_body do
    next "" unless it.rich_body.body.present?

    BlogPostRichTextContent.to_editor_html(it.rich_body.body.to_html)
  end
end
