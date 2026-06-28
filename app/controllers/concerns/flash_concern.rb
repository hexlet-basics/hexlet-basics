# typed: strict
# frozen_string_literal: true

module FlashConcern
  extend T::Sig
  extend T::Helpers
  requires_ancestor { ApplicationController }

  sig { params(key: T.untyped, options: T.untyped).returns(T.untyped) }
  def f(key, options = {})
    controller = self.class
    values = options[:values] || {}
    errors = options[:errors]

    msg = translate(key, options[:scope], controller, params[:action], values, errors)

    # NOTE color logging
    Rails.logger.debug(Term::ANSIColor.green("flash: #{msg}"))
    type = options[:type] || key
    if options[:now]
      flash.now[type] = msg
    else
      flash[type] = msg
    end
  end

  private

  sig { params(key: T.untyped, scope: T.untyped, controller: T.untyped, action: T.untyped, values: T.untyped, errors: T.untyped).returns(T.untyped) }
  def translate(key, scope, controller, action, values, errors = nil)
    keys = []
    lookup_controller = controller
    lookup_action = action

    if scope
      lookup_key =  scope.split("/")
      lookup_key << key

      keys << lookup_key.join(".").to_sym
    else
      while lookup_controller.superclass.name != "ActionController::Base"
        lookup_key = T.let([], T::Array[T.untyped])
        lookup_key << lookup_controller.controller_path.tr("/", ".")
        lookup_key << lookup_action
        lookup_key << key

        keys << lookup_key.join(".").to_sym

        lookup_controller = lookup_controller.superclass
        lookup_action = :base
      end
    end

    I18n.t(keys.shift, scope: :flash, default: keys, **values, errors: errors)
  end
end
