# frozen_string_literal: true

class ContainerBuilder
  def self.call(env)
    case env
    when 'production'
      ProductionContainer
    when 'development'
      DevelopmentContainer
    when 'test'
      TestContainer
    else
      raise "Container for #{env} does not exist"
    end
  end
end
