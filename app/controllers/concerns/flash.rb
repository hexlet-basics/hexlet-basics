# frozen_string_literal: true

module Flash
  def f(key, options = {})
    scope = :flash
    controller = self.class
    values = options[:values] || {}
    errors = options[:errors]

    keys = build_path(key, controller, params[:action])
    msg = I18n.t(keys.shift, scope: scope, default: keys, **values, errors: errors)

    Rails.logger.debug(Paint["flash: #{msg}", :green])
    type = options[:type] || key
    if options[:now]
      flash.now[type] = msg
    else
      flash[type] = msg
    end
  end

  private

  def build_path(type, controller, action)
    keys = []
    lookup_controller = controller
    lookup_action = action

    while lookup_controller.superclass.name != 'ActionController::Base'
      lookup_key = []
      lookup_key << lookup_controller.controller_path.tr('/', '.')
      lookup_key << lookup_action
      lookup_key << type

      keys << lookup_key.join('.').to_sym

      lookup_controller = lookup_controller.superclass
      lookup_action = :base
    end

    keys
  end
end
