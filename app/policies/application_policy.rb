# typed: strict
# frozen_string_literal: true

class ApplicationPolicy
  extend T::Sig

  sig { returns(T.untyped) }
  attr_reader :user, :record

  sig { params(user: T.untyped, record: T.untyped).void }
  def initialize(user, record)
    @user = T.let(user, T.untyped)
    @record = T.let(record, T.untyped)
  end

  sig { returns(T::Boolean) }
  def index?
    false
  end

  sig { returns(T::Boolean) }
  def show?
    false
  end

  sig { returns(T::Boolean) }
  def create?
    false
  end

  sig { returns(T::Boolean) }
  def new?
    create?
  end

  sig { returns(T::Boolean) }
  def update?
    false
  end

  sig { returns(T::Boolean) }
  def edit?
    update?
  end

  sig { returns(T::Boolean) }
  def destroy?
    false
  end

  class Scope
    extend T::Sig

    sig { params(user: T.untyped, scope: T.untyped).void }
    def initialize(user, scope)
      @user = T.let(user, T.untyped)
      @scope = T.let(scope, T.untyped)
    end

    sig { returns(T.untyped) }
    def resolve
      raise NoMethodError, "You must define #resolve in #{self.class}"
    end

    private

    sig { returns(T.untyped) }
    attr_reader :user, :scope
  end
end
