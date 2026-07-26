# typed: strict
# frozen_string_literal: true

require "sorbet-schema"

module ApplicationParamsStruct
  extend T::Sig

  module Base
    extend ActiveSupport::Concern
    extend T::Sig
    extend T::Helpers

    include ActiveModel::Validations
    mixes_in_class_methods ActiveModel::Validations::ClassMethods

    Attributes = T.type_alias { T::Hash[Symbol, T.untyped] }

    sig { returns(Attributes) }
    def attributes
      @attributes ||= T.let({}, T.nilable(Attributes))
    end
  end

  class InvalidParams < StandardError
    extend T::Sig

    sig { returns(ActiveModel::Errors) }
    attr_reader :errors

    sig { params(errors: ActiveModel::Errors).void }
    def initialize(errors)
      @errors = errors
      super("invalid struct params")
    end
  end

  sig do
    type_parameters(:S)
      .params(klass: T::Class[T.type_parameter(:S)], params: ActionController::Parameters)
      .returns(T.type_parameter(:S))
  end
  def self.from_params(klass, params)
    klass_t = T.unsafe(klass)
    props = klass_t.props
    permit_args = build_permit_args(props)
    permitted = params.permit(*T.unsafe(permit_args)).to_h.symbolize_keys
    props.each do |name, prop|
      type_object = prop[:type_object]
      unwrapped = T::Utils.unwrap_nilable(type_object) || type_object
      next unless unwrapped.is_a?(T::Types::TypedHash)
      raw = params[name]
      T.unsafe(permitted)[name] = raw.respond_to?(:to_unsafe_h) ? raw.to_unsafe_h : raw
    end
    cleaned = permitted.each_with_object({}) do |(name, value), hash|
      hash[name] = clean_value(value, props[name][:type_object])
    end
    result = klass_t.deserialize_from(:hash, cleaned)

    if result.success?
      struct = result.payload
      attach_attributes(struct, cleaned)
      return struct
    end

    fallback = klass_t.new
    fallback.errors.add(:base, result.error.message)
    Kernel.raise(InvalidParams.new(fallback.errors))
  end

  sig do
    type_parameters(:S)
      .params(klass: T::Class[T.type_parameter(:S)], hash: T::Hash[Symbol, T.untyped])
      .returns(T.type_parameter(:S))
  end
  def self.from_hash!(klass, hash)
    struct = T.unsafe(klass).new(**hash)
    attach_attributes(struct, hash)
    struct
  end

  sig { params(struct: Base, attributes: T::Hash[Symbol, T.untyped]).void }
  def self.attach_attributes(struct, attributes)
    T.unsafe(struct).instance_variable_set(:@attributes, attributes)
  end
  private_class_method :attach_attributes

  BOOLEAN_CASTER = ActiveModel::Type::Boolean.new
  private_constant :BOOLEAN_CASTER

  sig { params(value: T.untyped, type_object: T.untyped).returns(T.untyped) }
  def self.clean_value(value, type_object)
    if numeric_indexed_hash?(value)
      array = value.sort_by { |key, _| key.to_i }.map { |_, elem| clean_value(elem, nil) }
      return array.reject { |elem| elem.is_a?(String) && elem.empty? }
    end

    return nil if value.is_a?(String) && value.empty?
    return value.reject { |elem| elem.is_a?(String) && elem.empty? } if value.is_a?(Array)
    return BOOLEAN_CASTER.cast(value) if boolean_type?(type_object)

    value
  end
  private_class_method :clean_value

  sig { params(type_object: T.untyped).returns(T::Boolean) }
  def self.boolean_type?(type_object)
    unwrapped = T::Utils.unwrap_nilable(type_object) || type_object
    unwrapped == T::Boolean
  end
  private_class_method :boolean_type?

  sig { params(value: T.untyped).returns(T::Boolean) }
  def self.numeric_indexed_hash?(value)
    return false if !value.is_a?(Hash) || value.empty?

    value.keys.all? { |key| key.to_s.match?(/\A\d+\z/) }
  end
  private_class_method :numeric_indexed_hash?

  sig { params(props: T::Hash[Symbol, T::Hash[Symbol, T.untyped]]).returns(T::Array[T.untyped]) }
  def self.build_permit_args(props)
    scalars = []
    arrays = {}
    indexed_arrays = {}
    props.each do |name, prop|
      type_object = prop[:type_object]
      unwrapped = T::Utils.unwrap_nilable(type_object) || type_object
      if unwrapped.is_a?(T::Types::TypedArray)
        arrays[name] = []
        indexed_arrays[name] = {}
      elsif unwrapped.is_a?(T::Types::TypedHash)
        # extracted separately via to_unsafe_h in from_params
      else
        scalars << name
      end
    end
    arrays.empty? ? scalars : scalars + [ arrays, indexed_arrays ]
  end
  private_class_method :build_permit_args

  sig do
    type_parameters(:S)
      .params(klass: T::Class[T.type_parameter(:S)], params: ActionController::Parameters)
      .returns(T.type_parameter(:S))
  end
  def self.from_params!(klass, params)
    struct = T.unsafe(from_params(klass, params))
    return struct if struct.valid?

    Kernel.raise(InvalidParams.new(struct.errors))
  end
end
