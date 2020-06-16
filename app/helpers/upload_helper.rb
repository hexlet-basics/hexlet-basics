# frozen_string_literal: true

module UploadHelper
  def state_span_generator_for_upload(state)
    capitalized_state = state.capitalize

    case state
    when 'queued'
      content_tag(:span, capitalized_state, class: 'badge badge-secondary')
    when 'running'
      content_tag(:span, capitalized_state, class: 'badge badge-warning')
    when 'success'
      content_tag(:span, capitalized_state, class: 'badge badge-success')
    when 'fail'
      content_tag(:span, capitalized_state, class: 'badge badge-danger')
    end
  end
end
