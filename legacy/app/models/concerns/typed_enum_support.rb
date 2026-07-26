# typed: strict
# frozen_string_literal: true

module TypedEnumSupport
  extend T::Helpers
  extend ActiveSupport::Concern
  extend T::Sig

  requires_ancestor { ActiveRecord::Base }

  TypedEnumMetadata = T.type_alias do
    {
      enum_class: T.class_of(T::Enum),
      value_map: T::Hash[T.untyped, T.untyped]
    }
  end

  module ClassMethods
    extend T::Sig

    sig { returns(T::Hash[String, TypedEnumMetadata]) }
    def typed_enum_definitions
      @typed_enum_definitions ||= T.let({}, T.nilable(T::Hash[String, TypedEnumMetadata]))
    end

    sig { params(name: T.untyped, enum_class: T.untyped, options: T.untyped).returns(T.untyped) }
    def typed_enum(name, enum_class, **options)
      Kernel.raise ArgumentError, "#{enum_class} is not a T::Enum" unless enum_class.is_a?(Class) && enum_class < T::Enum

      value_map = normalize_typed_enum_values(enum_class)
      normalized_options = normalize_typed_enum_options(enum_class, options)

      typed_enum_definitions[name.to_s] = {
        enum_class: enum_class,
        value_map: value_map.stringify_keys
      }

      T.unsafe(self).enum(name, value_map, **normalized_options)
    end

    private

    sig do
      params(
        enum_class: T.class_of(T::Enum)
      ).returns(T::Hash[T.untyped, T.untyped])
    end
    def normalize_typed_enum_values(enum_class)
      enum_class.values.to_h do |enum_value|
        [ typed_enum_key(enum_class, enum_value), enum_value.serialize ]
      end
    end

    sig do
      params(
        enum_class: T.class_of(T::Enum),
        options: T::Hash[Symbol, T.untyped]
      ).returns(T::Hash[Symbol, T.untyped])
    end
    def normalize_typed_enum_options(enum_class, options)
      return options unless options.key?(:default)

      default_value = options[:default]
      return options.merge(default: nil) if default_value.nil?

      unless default_value.is_a?(enum_class)
        Kernel.raise ArgumentError, "default must be a #{enum_class.name}"
      end

      options.merge(default: default_value.serialize)
    end

    sig do
      params(
        enum_class: T.class_of(T::Enum),
        enum_value: T::Enum
      ).returns(String)
    end
    def typed_enum_key(enum_class, enum_value)
      serialized_value = enum_value.serialize
      # TODO: пока закоменировал так как надо ещё разобраться
      # 1. с генерацией типов для enums с помощью typelizer
      # так как он использует значения преобразованных значений ключей, что не соответствует реальным значениям данных
      # 2. также generated/enums
      #
      # if serialized_value.is_a?(String)
      #   return serialized_value if serialized_value.match?(/\A[a-z][a-z0-9_]*\z/)

      #   constant_name = enum_class.constants(false).find do |name|
      #     enum_class.const_get(name, false).equal?(enum_value)
      #   end

      #   Kernel.raise ArgumentError, "unable to resolve enum key for #{enum_value.inspect}" if constant_name.nil?

      #   return constant_name.to_s.underscore
      # end

      return serialized_value if serialized_value.is_a?(String)

      return serialized_value.to_s if serialized_value.is_a?(Symbol)

      # Перебираем константы enum-класса рефлексией, чтобы найти ключ по значению.
      # rubocop:disable Sorbet/ConstantsFromStrings
      constant_name = enum_class.constants(false).find do |name|
        enum_class.const_get(name, false).equal?(enum_value)
      end
      # rubocop:enable Sorbet/ConstantsFromStrings

      Kernel.raise ArgumentError, "unable to resolve enum key for #{enum_value.inspect}" if constant_name.nil?

      constant_name.to_s.underscore
    end
  end

  sig { params(enum_name: T.any(Symbol, String)).returns(T.nilable(T::Enum)) }
  def typed_enum_value(enum_name)
    definition = typed_enum_definition_for(enum_name)
    raw_value = public_send(enum_name)

    return if definition.nil? || raw_value.nil?

    serialized_value = definition[:value_map][raw_value.to_s]

    return if serialized_value.nil?

    definition[:enum_class].deserialize(serialized_value)
  end

  sig { params(enum_name: T.any(Symbol, String), value: T.nilable(T::Enum)).returns(T.untyped) }
  def typed_write_enum(enum_name, value)
    return public_send("#{enum_name}=", nil) if value.nil?

    definition = typed_enum_definition_for(enum_name)
    return public_send("#{enum_name}=", value.serialize) if definition.nil?

    key = definition[:value_map].key(value.serialize)

    public_send("#{enum_name}=", key || value.serialize)
  end

  private

  sig { params(enum_name: T.any(Symbol, String)).returns(T.nilable(TypedEnumMetadata)) }
  def typed_enum_definition_for(enum_name)
    T.unsafe(self.class).typed_enum_definitions[enum_name.to_s]
  end
end
