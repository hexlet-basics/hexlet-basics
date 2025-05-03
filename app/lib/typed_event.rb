# typed: strict

class TypedEvent < RailsEventStore::Event
  extend T::Sig

  sig { params(data: T::Hash[String, T.untyped], kwargs: T.untyped).void }
  def initialize(data = {}, **kwargs)
    # no need to transform keys
    super(data: data, **kwargs)
  end
  # sig { params(data: T::Hash[T.any(Symbol, String), T.untyped], kwargs: T.untyped).void }
  # def initialize(data:, **kwargs)
  #   symbolized_data = data.transform_keys(&:to_sym)
  #   super(data: symbolized_data, **kwargs)
  # end
end
