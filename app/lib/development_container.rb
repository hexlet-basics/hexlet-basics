# frozen_string_literal: true

class DevelopmentContainer
  extend Dry::Container::Mixin

  register 'docker_image_exercise_loader' do
    DockerImageExerciseLoader.new
  end
end
