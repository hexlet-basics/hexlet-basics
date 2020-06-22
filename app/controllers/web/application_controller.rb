# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  include AuthManagment
  include Flash
end
