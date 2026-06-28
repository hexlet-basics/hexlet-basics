# typed: strict

class RuSuffixConstraint
  extend T::Sig

  sig { params(request: T.untyped).returns(T::Boolean) }
  def self.matches?(request)
    request.params[:suffix] == "ru"
  end
end
