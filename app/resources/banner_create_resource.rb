# typed: strict

class BannerCreateResource < ApplicationResource
  typelize_from Banner

  attributes :id,
    :state,
    :background,
    :locale,
    body: [ String, true ],
    url: [ String, true ]

  typelize id: [ :number, nullable: true ]
  typelize starts_at: [ :string, nullable: true ]
  typelize finishes_at: [ :string, nullable: true ]

  attribute :starts_at do |banner|
    banner.starts_at&.iso8601
  end

  attribute :finishes_at do |banner|
    banner.finishes_at&.iso8601
  end
end
