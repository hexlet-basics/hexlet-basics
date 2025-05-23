require "faraday"
require "httpx/adapters/faraday"
require "json"

class N8nClient < N8nClientInterface
  def initialize(base_url: ENV.fetch("N8N_BASE_URL"))
    @conn = Faraday.new(url: base_url) do |f|
        f.request :json
        f.response :raise_error
        f.adapter :httpx
      end
  end

  def trigger_lead_created_workflow(serializer)
    response = @conn.post(ENV.fetch("N8N_LEAD_CREATED_PATH")) do |req|
      req.headers["Content-Type"] = "application/json"
      req.body = serializer.to_h
    end

    JSON.parse(response.body, symbolize_names: true)
  end
end
