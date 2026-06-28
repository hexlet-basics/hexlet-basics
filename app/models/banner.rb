# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: banners
#
#  id          :bigint           not null, primary key
#  background  :string           default("cta_gradient"), not null
#  body        :text             not null
#  finishes_at :datetime
#  locale      :string           not null
#  starts_at   :datetime
#  state       :string           default("draft"), not null
#  url         :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
# Indexes
#
#  index_banners_on_locale_and_state  (locale,state)
#
class Banner < ApplicationRecord
  class State < T::Enum
    enums do
      Draft = new("draft")
      Published = new("published")
      Archived = new("archived")
    end
  end

  class Locale < T::Enum
    enums do
      Ru = new("ru")
      En = new("en")
    end
  end

  class Background < T::Enum
    enums do
      CtaGradient = new("cta_gradient")
      Dark = new("dark")
      Blue = new("blue")
    end
  end

  include BannerRepository

  typed_enum :state, State, suffix: true, validate: true, default: State::Draft
  typed_enum :locale, Locale, suffix: true, validate: true
  typed_enum :background, Background, suffix: true, validate: true, default: Background::CtaGradient

  validates :body, presence: true
  validates :locale, presence: true

  sig { params(_auth_object: T.untyped).returns(T.untyped) }
  def self.ransackable_attributes(_auth_object = nil)
    [ "id", "created_at", "state", "locale", "starts_at", "finishes_at" ]
  end

  sig { params(_auth_object: T.untyped).returns(T.untyped) }
  def self.ransackable_associations(_auth_object = nil)
    []
  end

  sig { returns(String) }
  def to_s
    body.to_s.truncate(50)
  end
end
