module AhoyConcern
  extend ActiveSupport::Concern

  included do
    skip_before_action :track_ahoy_visit
    before_action :track_utm_visit_with_ahoy
  end

  private

  def track_utm_visit_with_ahoy
    return unless params[:utm_source].present?

    # NOTE write Ahoy visitor_token & visit_token cookies as we skip automatic methods
    ahoy.send(:set_visitor_cookie)
    ahoy.send(:set_visit_cookie)

    ahoy.track_visit
  end
end
