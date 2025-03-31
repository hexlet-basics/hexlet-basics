# frozen_string_literal: true

# == Schema Information
#
# Table name: users
#
#  id                   :integer          not null, primary key
#  admin                :boolean
#  confirmation_token   :string(255)
#  email                :string(255)
#  email_delivery_state :string(255)
#  facebook_uid         :string(255)
#  first_name           :string(255)
#  github_uid           :integer
#  help                 :boolean
#  last_name            :string(255)
#  locale               :string(255)
#  nickname             :string(255)
#  password_digest      :string(255)
#  reset_password_token :string(255)
#  state                :string(255)
#  created_at           :datetime         not null
#  updated_at           :datetime         not null
#
# Indexes
#
#  index_users_on_email  (email) UNIQUE
#
class User < ApplicationRecord
  include UserRepository
  include AASM

  has_secure_password validations: false

  def self.ransackable_attributes(_auth_object = nil)
    %w[email first_name last_name]
  end

  def self.ransackable_associations(_auth_object = nil)
    [ "language_members" ]
  end

  validates :email, presence: true,
                    uniqueness: { case_sensitive: false },
                    'valid_email_2/email': { mx: true },
                    unless: :removed?

  has_many :lesson_members, class_name: "Language::Lesson::Member", dependent: :destroy
  has_many :lessons, through: :lesson_members, class_name: "Language::Lesson"
  has_many :language_members, class_name: "Language::Member", dependent: :destroy
  has_many :accounts, dependent: :destroy
  has_many :reviews, dependent: :destroy
  has_many :blog_posts, dependent: :destroy

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
