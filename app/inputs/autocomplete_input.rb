# frozen_string_literal: true

class AutocompleteInput < SimpleForm::Inputs::CollectionSelectInput
  def input(wrapper_options = {})
    # raise self.reflection.name.inspect
    options[:collection] = [object.send(:"#{reflection.name}")]
    template.append_javascript_packs('autocomplete')
    # template.append_stylesheet_packs('autocomplete')
    # https://select2.org/configuration/defaults
    # https://stackoverflow.com/questions/48007094/select2-not-using-data-attributes
    input_html_options['data-ajax--url'] = options[:source]
    # input_html_options['data-ajax--cache'] = true

    super
  end

  def input_html_classes
    super.push('select2')
  end
end
