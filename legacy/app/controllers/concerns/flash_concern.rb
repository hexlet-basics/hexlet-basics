# typed: strict
# frozen_string_literal: true

module FlashConcern
  extend T::Sig
  extend T::Helpers
  requires_ancestor { ApplicationController }

  sig do
    params(
      key: T.any(String, Symbol),
      scope: T.nilable(String),
      values: T::Hash[Symbol, T.untyped],
      errors: T.untyped,
      type: T.nilable(T.any(String, Symbol)),
      now: T::Boolean
    ).void
  end
  def f(key, scope: nil, values: {}, errors: nil, type: nil, now: false)
    controller = T.cast(self.class, T.class_of(ActionController::Base))

    msg = translate(key, scope, controller, params[:action], values, errors)

    # NOTE color logging
    Rails.logger.debug(Term::ANSIColor.green("flash: #{msg}"))
    flash_type = type || key
    if now
      flash.now[flash_type] = msg
    else
      flash[flash_type] = msg
    end
  end

  private

  sig do
    params(
      key: T.any(String, Symbol),
      scope: T.nilable(String),
      controller: T.class_of(ActionController::Base),
      action: T.nilable(String),
      values: T::Hash[Symbol, T.untyped],
      errors: T.untyped
    ).returns(String)
  end
  def translate(key, scope, controller, action, values, errors = nil)
    keys = T.let([], T::Array[Symbol])
    lookup_controller = T.let(controller, T.nilable(T.class_of(ActionController::Base)))
    lookup_action = T.let(action, T.nilable(T.any(String, Symbol)))

    if scope
      keys << [ *scope.split("/"), key ].join(".").to_sym
    else
      while lookup_controller && lookup_controller.superclass&.name != "ActionController::Base"
        path = lookup_controller.controller_path.tr("/", ".")
        keys << [ path, lookup_action, key ].join(".").to_sym

        lookup_controller = T.cast(lookup_controller.superclass, T.nilable(T.class_of(ActionController::Base)))
        lookup_action = :base
      end
    end

    I18n.t(T.must(keys.shift), scope: :flash, default: keys, **values, errors: errors)
  end
end
