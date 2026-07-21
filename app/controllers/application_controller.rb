# typed: strict

class ApplicationController < ActionController::Base
  extend T::Sig

  # meta-tags includes this into ActionController::Base via its Railtie on_load
  # hook (which tapioca can't see); state it explicitly so set_meta_tags resolves
  # for typed controllers without an RBI shim. Idempotent at runtime.
  include MetaTags::ControllerHelper
  # browser gem also mixes this in via its Railtie on_load hook (invisible to
  # tapioca); state it explicitly so `browser` resolves on controllers. Idempotent.
  include Browser::ActionController

  include ActiveStorage::SetCurrent
  include Pundit::Authorization
  include Pagy::Method
  include Authentication
  include LocaleConcern
  include EventConcern

  # rescue_from Pagy::OverflowError, with: :redirect_to_last_page
  # rescue_from Pagy::VariableError, with: :redirect_to_last_page

  sig { returns(T::Hash[Symbol, T.untyped]) }
  def default_url_options
    { suffix: params[:suffix] }
  end

  # def redirect_to_last_page(exception)
  #   redirect_to url_for(page: exception.pagy.last), notice: "Page ##{params[:page]} is overflowing. Showing page #{exception.pagy.last} instead."
  # end

  private

  # Typed wrapper over Pagy::Method#pagy. The gem RBI types `pagy` as
  # `(paginator = :offset, collection, **options)`, so Sorbet rejects a
  # single-argument call; pass the default `:offset` paginator explicitly and
  # expose a real `[Pagy, records]` return so call sites stay fully typed.
  sig { params(collection: T.untyped, options: T.untyped).returns([ Pagy::Offset, T.untyped ]) }
  def paginate(collection, **options)
    pagy(:offset, collection, **options)
  end
end
