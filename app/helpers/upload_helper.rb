# frozen_string_literal: true

module UploadHelper
  def uploaded_state_class(state)
    case state
    when 'not_run'
      'badge badge-secondary'
    when 'running'
      'badge badge-warning'
    when 'success'
      'badge badge-success'
    when 'failed'
      'badge badge-danger'
    end
  end
end
