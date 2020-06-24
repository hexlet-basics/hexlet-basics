# frozen_string_literal: true

class TestContainer
  extend Dry::Container::Mixin

  register 'docker_image_exercise_loader' do
    :success
  end
end
