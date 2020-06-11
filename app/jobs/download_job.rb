# frozen_string_literal: true

class DownloadJob < ApplicationJob
  def perform(exercise)
    system("make web-exercises-load-#{exercise}")
  end
end
