module AhoyConcern
  extend ActiveSupport::Concern

  included do
    before_action :track_utm_visit_with_ahoy
  end

  private

  def track_utm_visit_with_ahoy
    if params[:utm_source].present?
      ahoy.track_visit
    end
  end
end
