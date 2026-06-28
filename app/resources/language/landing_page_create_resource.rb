# typed: strict

class Language::LandingPageCreateResource < ApplicationResource
  typelize_from Language::LandingPage

  attributes :state,
    :footer,
    :listed,
    :main,
    :outcomes_image,
    :id,
    :language_id,
    :landing_page_to_redirect_id,
    slug: [ String, true ],
    footer_name: [ String, true ],
    name: [ String, true ],
    order: [ String, true ],
    meta_title: [ String, true ],
    meta_description: [ String, true ],
    header: [ String, true ],
    description: [ String, true ],
    used_in_header: [ String, true ],
    used_in_description: [ String, true ],
    outcomes_header: [ String, true ],
    outcomes_description: [ String, true ]

  typelize outcomes_image: "File | null"
  typelize id: [ :number, nullable: true ]
  typelize language_id: [ :number, nullable: true ]
  typelize landing_page_to_redirect_id: [ :number, nullable: true ]
end
