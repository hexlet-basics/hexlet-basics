# frozen_string_literal: true

class Web::ForSchoolTeachersController < Web::ApplicationController
  def index
    render inertia: true, props: {}
  end
end
