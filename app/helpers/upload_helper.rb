# frozen_string_literal: true

module UploadHelper
  def uploaded_state_class(state)
    case state
    when 'created'
      'badge badge-secondary'
    when 'building'
      'badge badge-warning'
    when 'built'
      'badge badge-success'
    end
  end
end
