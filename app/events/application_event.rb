# typed: strict

class ApplicationEvent < RailsEventStore::Event
  extend T::Sig
  extend T::Helpers
  abstract!

  sig { params(data: T::Hash[String, T.untyped], kwargs: T.untyped).void }
  def initialize(data = {}, **kwargs)
    # no need to transform keys
    super(data: data, **kwargs)
  end

  EventMetadata = T.type_alias { T::Hash[Symbol, T.untyped] }

  sig { returns(EventMetadata) }
  def as_json
    {
      type: event_type,
      name: T.cast(self.class.const_get(:NAME), String),
      event_id: event_id,
      metadata: metadata.to_h,
      data: data
    }
  end
end
