# frozen_string_literal: true

require 'pty'

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

      yield $CHILD_STATUS.inspect if block_given? && !success?($CHILD_STATUS)

      force_fail ? false : success?($CHILD_STATUS)
    end

    def success?(child_status)
      child_status&.exitstatus&.zero?
    end
  end
end
