# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: users
#
#  id                       :bigint           not null, primary key
#  admin                    :boolean
#  assistant_messages_count :integer          default(0)
#  confirmation_token       :string(255)
#  contact_method           :string
#  contact_value            :string
#  email                    :string(255)
#  email_delivery_state     :string(255)
#  facebook_uid             :string(255)
#  first_name               :string(255)
#  github_uid               :integer
#  help                     :boolean
#  last_name                :string(255)
#  locale                   :string(255)
#  nickname                 :string(255)
#  password_digest          :string(255)
#  phone                    :string
#  phone_verified_at        :datetime
#  state                    :string(255)
#  created_at               :datetime         not null
#  updated_at               :datetime         not null
#  webauthn_id              :string
#
# Indexes
#
#  index_users_on_LOWER_email  (lower((email)::text)) UNIQUE
#  index_users_on_email        (email) UNIQUE
#  index_users_on_phone        (phone) UNIQUE WHERE (phone IS NOT NULL)
#  index_users_on_webauthn_id  (webauthn_id) UNIQUE WHERE (webauthn_id IS NOT NULL)
#
class User < ApplicationRecord
  class ContactMethod < T::Enum
    enums do
      Telegram = new("telegram")
      Phone = new("phone")
      Whatsapp = new("whatsapp")
    end
  end

  include UserRepository
  include AASM

  acts_as_taggable_on :tags

  has_secure_password validations: false, reset_token: { expires_in: 15.minutes }

  # Passwordless email sign-in. Bound to email so the link dies if the email changes.
  generates_token_for :magic_link, expires_in: 15.minutes do
    T.bind(self, User)
    email
  end

  normalizes :email, with: ->(email) { email.strip.downcase.presence }
  normalizes :phone, with: ->(phone) { Phonelib.parse(phone).e164.presence }

  sig { params(_auth_object: T.untyped).returns(T.untyped) }
  def self.ransackable_attributes(_auth_object = nil)
    %w[id email first_name last_name created_at]
  end

  sig { params(_auth_object: T.untyped).returns(T.untyped) }
  def self.ransackable_associations(_auth_object = nil)
    [ "language_members" ]
  end

  # Email is optional: users can sign in by phone or a linked social account.
  # Postgres unique indexes treat multiple NULLs as distinct, so blank email is normalized to nil.
  validates :email,
    uniqueness: { case_sensitive: false },
    'valid_email_2/email': { mx: true },
    allow_blank: true,
    unless: :removed?

  validates :phone, uniqueness: true, allow_blank: true, unless: :removed?

  validate :identifier_present, unless: :removed?

  validates :first_name, length: { maximum: 40 },
    format: { with: UsefulRegexp.without_spec_chars },
    allow_blank: true

  validates :last_name, length: { maximum: 40 },
    format: { with: UsefulRegexp.without_spec_chars },
    allow_blank: true

  validates :nickname, length: { maximum: 40 },
    format: { with: UsefulRegexp.without_spec_chars },
    allow_blank: true

  has_many :visits, class_name: "Ahoy::Visit"
  has_many :events, through: :visits, class_name: "Ahoy::Event"

  has_many :lesson_members, class_name: "Language::Lesson::Member", dependent: :destroy
  has_many :lessons, through: :lesson_members, class_name: "Language::Lesson"
  has_many :language_members, class_name: "Language::Member", dependent: :destroy
  has_many :sessions, dependent: :destroy
  has_many :accounts, dependent: :destroy
  has_many :credentials, dependent: :destroy
  has_many :reviews, dependent: :destroy
  has_many :blog_posts, foreign_key: "creator_id", dependent: :destroy

  has_many :survey_answers, class_name: "Survey::Answer", dependent: :destroy
  has_many :survey_answers_surveys, through: :survey_answers, source: :survey
  has_many :survey_answers_survey_items, through: :survey_answers, source: :survey_item
  has_many :survey_scenario_members, class_name: "Survey::Scenario::Member", dependent: :destroy
  has_many :survey_scenarios, through: :survey_scenario_members, source: :scenario

  has_many :assistant_messages, class_name: "Language::Lesson::Member::Message"

  has_one :book_request
  has_many :leads

  typed_enum :contact_method, ContactMethod, suffix: true

  aasm :state do
    state :active, initial: true
    state :waiting_confirmation
    state :removed

    event :activate do
      transitions to: :active
    end

    event :mark_as_removed do
      transitions to: :removed

      after do
        clean_fields
        remove_accounts
      end
    end
  end


  sig { returns(String) }
  def to_s
    if first_name? || last_name?
      return "#{first_name} #{last_name}"
    end
    email.presence || phone.presence || "User ##{id}"
  end

  sig { returns(T::Boolean) }
  def should_be_lead?
    tag_list.include?("should_be_lead")
  end

  private

  # A user must be reachable/identifiable by at least one credential:
  # email, phone, or a linked social account (built in memory or persisted).
  sig { void }
  def identifier_present
    return if email? || phone? || accounts.any?

    errors.add(:base, :identifier_required)
  end

  sig { returns(T.untyped) }
  def clean_fields
    fields = %w[
      first_name
      last_name
      nickname
      password_digest
      confirmation_token
      email
      phone
    ]

    fields.each do |field|
      send :"#{field}=", nil
    end

    save!
  end

  sig { returns(T.untyped) }
  def remove_accounts
    accounts.clear
  end
end
