# frozen_string_literal: true

class EventsTypesExporter
  def initialize(output_path: Rails.root.join("app/javascript/types/events.ts"))
    @output_path = output_path
  end

  def call
    events = build_events
    FileUtils.mkdir_p(@output_path.dirname)
    File.write(@output_path, render(events))
    puts "✅ Exported #{events.size} event types to #{@output_path}"
  end

  private

  def build_events
    ApplicationEvent.descendants
      .sort_by(&:name)
      .map do |klass|
        event_type = klass.name
        data_fields = data_shape_fields(klass)

        {
          class_name: klass.name.gsub("::", ""),
          event_type: event_type,
          event_name: klass::NAME,
          data_fields: data_fields
        }
      end
      .sort_by { |event| event[:event_type] }
  end

  def data_shape_fields(klass)
    data_shape = klass.const_get(:DataShape)
    type = unwrap_type_alias(data_shape)
    return [] unless type.is_a?(T::Types::FixedHash)

    type.types.map do |key, value|
      [ key, ts_type_for(value) ]
    end
  rescue NameError
    []
  end

  def unwrap_type_alias(type)
    return type.aliased_type if type.respond_to?(:aliased_type)

    type
  end

  def ts_type_for(type)
    type = unwrap_type_alias(type)

    case type
    when T::Types::TypedClass
      ts_typed_class_type_for(type)
    when T::Types::Simple
      ts_simple_type_for(type)
    when T::Types::TypedArray
      "#{ts_type_for(type.type)}[]"
    when T::Types::TypedHash
      key_type = ts_type_for(type.keys)
      key_type = "string" if key_type == "unknown"
      value_type = ts_type_for(type.values)
      "Record<#{key_type}, #{value_type}>"
    when T::Types::Union, T::Private::Types::SimplePairUnion
      ts_union_type_for(type)
    when T::Types::Untyped
      "unknown"
    when T::Types::FixedHash
      "Record<string, unknown>"
    else
      "unknown"
    end
  end

  def ts_typed_class_type_for(type)
    raw_type = type.raw_type

    if raw_type < T::Enum
      ts_enum_type_for(raw_type)
    else
      "unknown"
    end
  end

  def ts_enum_type_for(enum_class)
    values = enum_class.values.map { |value| value.serialize.to_s }.uniq
    return "string" if values.empty?

    values.map(&:inspect).join(" | ")
  end

  def ts_union_type_for(type)
    union_types = type.types.map { |inner| ts_type_for(inner) }.uniq
    return "unknown" if union_types.empty?

    union_types.join(" | ")
  end

  def ts_simple_type_for(type)
    raw_type = type.raw_type

    if raw_type < T::Enum
      return ts_enum_type_for(raw_type)
    end

    case raw_type.name
    when "String", "Symbol"
      "string"
    when "Integer", "Float"
      "number"
    when "TrueClass", "FalseClass"
      "boolean"
    when "Time", "Date", "DateTime"
      "string"
    when "NilClass"
      "undefined"
    else
      "unknown"
    end
  end

  def render(events)
    lines = []
    lines << "// This file is auto-generated. Do not edit manually."
    lines << "// Run `bin/rails app:export_events_to_ts` to update."
    lines << ""

    events.each do |event|
      lines << "export interface #{event[:class_name]} {"
      lines << "  type: \"#{event[:event_type]}\";"
      lines << "  name: \"#{event[:event_name]}\";"
      lines << "  event_id: string;"
      lines << "  metadata: { name: string };"

      if event[:data_fields].empty?
        lines << "  data: Record<string, never>;"
      else
        lines << "  data: {"
        event[:data_fields].each do |key, type|
          lines << "    #{key}: #{type};"
        end
        lines << "  };"
      end

      lines << "}"
      lines << ""
    end

    lines << "export type ApplicationEvent ="
    events.each do |event|
      lines << "  | #{event[:class_name]}"
    end
    "#{lines.join("\n")};"
  end
end
