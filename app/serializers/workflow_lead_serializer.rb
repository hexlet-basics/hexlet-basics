class WorkflowLeadSerializer
  # Соответствие ключей входных данных и ID полей в AmoCRM
  LEAD_FIELD_MAPPER = {}.freeze

  CONTACT_FIELD_MAPPER = {
    user_id:   1_030_881,
    email:     { field_id: 316_907, enum_id: 158_463 }, # Личная почта
    phone:     { field_id: 316_905, enum_id: 158_453 }, # Мобильный телефон
    telegram:  925_081,
    whatsapp:  1_012_781
  }.freeze

  DEFAULT_TAG = { name: "cb" }.freeze

  def initialize(event)
    @data = event.data
  end

  def payload
    lead = base_payload
    lead[:custom_fields] = build_custom_fields(LEAD_FIELD_MAPPER)

    embedded = {
      contacts: [ contact_payload ],
      tags:     [ DEFAULT_TAG ]
    }

    lead[:_embedded] = embedded

    {
      lead: [ lead ],
      note: notes_payload
    }
  end

  private

  def base_payload
    {
      name:                @data[:email],
      pipeline_id:         5_101_747,
      status_id:           45_861_256,
      responsible_user_id: 7_877_026
    }
  end

  def notes_payload
    text = @data[:survey_answers_data]
             .map { |qa| "Вопрос: #{qa[:question]}\nОтвет: #{qa[:answer]}" }
             .join("\n\n")

    { text: text }
  end

  def contact_payload
    {
      name:                 @data[:user_name],
      custom_fields_values: build_custom_fields(CONTACT_FIELD_MAPPER)
    }
  end

  def build_custom_fields(mapper)
    mapper.each_with_object([]) do |(key, config), arr|
      next unless @data[key]&.present?

      field_id, enum_id = if config.is_a?(Hash)
                            [ config[:field_id], config[:enum_id] ]
      else
                            [ config, nil ]
      end

      value_hash = { value: @data[key] }
      value_hash[:enum_id] = enum_id if enum_id

      arr << {
        field_id: field_id,
        values:   [ value_hash ]
      }
    end
  end
end
