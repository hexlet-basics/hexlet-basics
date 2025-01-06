# typed: true

# DO NOT EDIT MANUALLY
# This file was pulled from a central RBI files repository.
# Please run `bin/tapioca annotations` to update it.

module AASM
  mixes_in_class_methods AASM::ClassMethods

  sig { params(name: T.any(Symbol, String)).returns(AASM::InstanceBase) }
  def aasm(name = :default); end
end

module AASM::ClassMethods
  sig { params(args: T.untyped, block: T.nilable(T.proc.void)).void }
  def aasm(*args, &block); end
end

class AASM::Base
  sig { params(callbacks: Symbol, block: T.nilable(T.proc.bind(T.untyped).void)).void }
  def after_all_transitions(*callbacks, &block); end

  sig { params(callbacks: Symbol, block: T.nilable(T.proc.bind(T.untyped).void)).void }
  def after_all_transactions(*callbacks, &block); end

  sig { params(callbacks: Symbol, block: T.nilable(T.proc.bind(T.untyped).void)).void }
  def before_all_transactions(*callbacks, &block); end

  sig { params(callbacks: Symbol, block: T.nilable(T.proc.bind(T.untyped).void)).void }
  def before_all_events(*callbacks, &block); end
end

class AASM::Core::Event
  # @shim: added through `instance_eval` by AASM::DslHelper::Proxy
  sig { params(name: Symbol, block: T.proc.bind(T.untyped).void).void }
  def before(*name, &block); end

  # @shim: added through `instance_eval` by AASM::DslHelper::Proxy
  sig { params(name: Symbol, block: T.proc.bind(T.untyped).void).void }
  def after(*name, &block); end
end
