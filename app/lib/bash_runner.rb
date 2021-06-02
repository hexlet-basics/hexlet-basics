# frozen_string_literal: true

require 'pty'
# rubocop:disable Style/SpecialGlobalVars
class BashRunner
  class << self
    def start(command, force_fail: false)
      PTY.spawn(command, $stderr => $stdout) do |r, _w, pid|
        begin
          r.each { |line| yield line if block_given? }
        rescue Errno::EIO => e
          Rails.logger.warn e
        end
        Process.wait(pid)
      end

      # NOTE Test child status
      status = $?

      yield status.inspect if block_given? && !success?(status)

      force_fail ? false : status
    end

    def success?(child_status)
      child_status&.exitstatus&.zero?
    end
  end
end
# rubocop:enable Style/SpecialGlobalVars
