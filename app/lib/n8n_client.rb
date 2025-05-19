require "faraday"
require "httpx/adapters/faraday"

class N8nClient
  def initialize(base_url: ENV.fetch("N8N_BASE_URL"))
    @conn = Faraday.new(url: base_url) do |f|
      f.request :json
      f.response :raise_error
      f.adapter :httpx
    end
  end

  def send(payload)
    response = @conn.post(
      ENV.fetch("N8N_WORKFLOW_PATH"),
      payload
    )

    JSON.parse(response.body, symbolize_names: true)
  end
end
