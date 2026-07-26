# typed: strict
# frozen_string_literal: true

class Web::Admin::Management::ApplicationController < Web::Admin::ApplicationController
  # Управление пользователями, ролями и staff-участниками — строго для суперпользователя (admin),
  # чтобы staff не мог выдать права самому себе.
  before_action :authenticate_admin!
end
