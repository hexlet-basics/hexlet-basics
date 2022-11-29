# frozen_string_literal: true

class AasmStateEventInput < SimpleForm::Inputs::CollectionSelectInput
  def collection
    object.aasm(state_machine_props.name).events_for_select
  end

  def input(wrapper_options)
    label_method = :default_display_name
    namespace = state_machine_props.config.namespace
    value_method = if namespace
                     ->(event) { "#{event.name}_#{namespace}" }
                   else
                     :name
                   end

    current_state_tag = template.tag.span("#{I18n.t('.inputs.current_state')} #{state_machine.human_state}", class: 'text-muted small')

    merged_html_options = merge_wrapper_options(input_html_options, wrapper_options)
    merged_input_options = input_options.merge(include_blank: 'select action')

    out = @builder.collection_select(
      attribute_name, state_machine.events, value_method, label_method,
      merged_input_options, merged_html_options
    )
    (out << current_state_tag).html_safe
  end

  private

  def clean_attribute_name
    attribute_name.to_s.gsub('_event', '')
  end

  def state_machine_props
    return @state_machine_props if @state_machine_props

    aasm_store = AASM::StateMachineStore.fetch(object.class, true)
    aasm_store_items = aasm_store.machine_names.map { |name| aasm_store.machine(name) }
    @state_machine_props = aasm_store_items.find do |sm|
      column = if sm.config.namespace
                 "#{sm.config.column}_#{sm.config.namespace}"
               else
                 sm.config.column.to_s
               end

      column == clean_attribute_name
    end

    unless @state_machine_props
      raise ArgumentError, "AASM state machine for field <#{clean_attribute_name}> not found in class #{object.class}"
    end

    @state_machine_props
  end

  def state_machine
    object.aasm(@state_machine_props.name)
  end
end
