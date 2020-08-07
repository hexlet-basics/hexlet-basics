# frozen_string_literal: true

class StateEventInput < SimpleForm::Inputs::CollectionSelectInput
  def collection
    options[:collection].map do |event|
      [I18n.t(event.name, scope: ['events', object.class.name.underscore]), event.name]
    end
  end

  def input(wrapper_options)
    label_method = :first
    value_method = :second

    current_state_name = object.aasm(clean_attribute_name).human_state
    current_state = template.tag.span("#{I18n.t('.inputs.current_state')} #{current_state_name}", class: 'mx-1')

    merged_html_options = merge_wrapper_options(input_html_options, wrapper_options)
    merged_input_options = input_options.merge(include_blank: 'select action')

    out = @builder.collection_select(
      attribute_name, collection, value_method, label_method,
      merged_input_options, merged_html_options
    )
    (out << current_state).html_safe # rubocop:disable Rails/OutputSafety
  end

  def clean_attribute_name
    attribute_name.to_s.gsub '_event', ''
  end
end
