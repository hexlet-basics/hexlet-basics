# frozen_string_literal: true

module AASM
  class InstanceBase
    def events_for_select
      events(permitted: true).map do |event|
        [@instance.class.aasm(@name).human_event_name(event), event.name]
      end
    end
  end

  module ExtBase
    def initialize(klass, name, state_machine, options = {}, &)
      super

      # NOTE: AASM по дефолту не рейзит исключение если при вызове события с ! зафейлились валидации в ActiveRecord
      @state_machine.config.whiny_persistence = true
      setup_event_direct_assignment(@name)
    end

    def events_for_select
      events.map do |event|
        [human_event_name(event), event.name]
      end
    end

    private

    def setup_event_direct_assignment(aasm_name)
      column = @state_machine.config.column
      namespace = @state_machine.config.namespace
      attribute_name = if namespace
                         "#{column}_#{namespace}_event"
                       else
                         "#{column}_event"
                       end

      @klass.attribute attribute_name, :string
      @klass.before_save :"update_state_by_event_attribute_#{attribute_name}"

      @klass.send(:define_method, "update_state_by_event_attribute_#{attribute_name}") do
        event = send(attribute_name)&.to_sym
        if event.present?
          aasm(aasm_name).fire(event)
          send(:"#{attribute_name}=", nil)
        end
      end
    end
  end

  class Base
    prepend ExtBase
  end
end
