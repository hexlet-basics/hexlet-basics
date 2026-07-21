# typed: true

# `ActiveStorage::Attached::One` reaches `#variant`, `#representation` and
# `#preview` through `delegate_missing_to :attachment`, so tapioca can't see
# them and Sorbet reports them as missing. Declare the delegated signatures
# here (mirrors `ActiveStorage::Attachment`), so call sites stay fully typed
# without `T.unsafe`. Shims survive `tapioca gem`/`tapioca dsl` regeneration.
class ActiveStorage::Attached::One
  sig do
    params(transformations: T.any(Symbol, T::Hash[Symbol, T.untyped]))
      .returns(T.any(ActiveStorage::Variant, ActiveStorage::VariantWithRecord))
  end
  def variant(transformations); end

  sig do
    params(transformations: T.any(Symbol, T::Hash[Symbol, T.untyped]))
      .returns(T.any(ActiveStorage::Variant, ActiveStorage::VariantWithRecord, ActiveStorage::Preview))
  end
  def representation(transformations); end

  sig { params(transformations: T::Hash[Symbol, T.untyped]).returns(ActiveStorage::Preview) }
  def preview(transformations); end
end
