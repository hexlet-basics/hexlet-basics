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
#  reset_password_token     :string(255)
#  state                    :string(255)
#  created_at               :datetime         not null
#  updated_at               :datetime         not null
#
# Indexes
#
#  index_users_on_email  (email) UNIQUE
#
class User < ApplicationRecord
  include UserRepository
  include AASM

  acts_as_taggable_on :tags

  has_secure_password validations: false

  def self.ransackable_attributes(_auth_object = nil)
    %w[id email first_name last_name created_at]
  end

  def self.ransackable_associations(_auth_object = nil)
    [ "language_members" ]
  end

  validates :email, presence: true,
                    uniqueness: { case_sensitive: false },
                    'valid_email_2/email': { mx: true },
                    unless: :removed?

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
  has_many :accounts, dependent: :destroy
  has_many :reviews, dependent: :destroy
  has_many :blog_posts, dependent: :destroy

  has_many :survey_answers, class_name: "Survey::Answer"
  has_many :survey_answers_surveys, through: :survey_answers, source: :survey
  has_many :survey_answers_survey_items, through: :survey_answers, source: :survey_item
  has_many :survey_scenario_members, class_name: "Survey::Scenario::Member"
  has_many :survey_scenarios, through: :survey_scenario_members, source: :scenario

  has_many :assistant_messages, class_name: "Language::Lesson::Member::Message"

  has_one :book_request
  has_many :leads

  enum :contact_method, { telegram: "telegram", phone: "phone", whatsapp: "whatsapp" }, suffix: true

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

  def initialize(attrs = nil)
    defaults = {
      locale: I18n.locale
    }
    attrs_with_defaults = attrs ? defaults.merge(attrs) : defaults
    super(attrs_with_defaults)
  end

  def guest?
    false
  end

  def valid_password?(password)
    return false if password_digest.nil?

    authenticate(password)
  end

  def to_s
    if first_name? || last_name?
      return "#{first_name} #{last_name}"
    end

    email
  end

  def should_be_lead?
    tag_list.include?("should_be_lead")
  end

  private

  def clean_fields
    fields = %w[
      first_name
      last_name
      nickname
      password_digest
      reset_password_token
      confirmation_token
      email
    ]

    fields.each do |field|
      send :"#{field}=", nil
    end

    save!
  end

  def remove_accounts
    accounts.clear
  end
end
