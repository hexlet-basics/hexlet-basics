# frozen_string_literal: true

class User::Account < ApplicationRecord
  belongs_to :user
end
