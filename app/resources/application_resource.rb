# typed: strict

class ApplicationResource
  extend T::Sig
  include Alba::Resource
  include Typelizer::DSL

  meta do
    {}
  end

  sig { returns(T.untyped) }
  def to_model
    object
  end
end
