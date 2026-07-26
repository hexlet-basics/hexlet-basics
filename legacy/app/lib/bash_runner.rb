# typed: strict

require "English" # NOTE Без явного указания $CHILD_STATUS nil
require "pty"

class BashRunner
  class << self
    extend T::Sig

    sig { params(command: T.untyped, force_fail: T::Boolean, block: T.untyped).returns(T.untyped) }
    def start(command, force_fail: false, &block)
      PTY.spawn(command, $stderr => $stdout) do |r, _w, pid|
        begin
          r.each { |line| yield line if block_given? }
        rescue Errno::EIO => e
          Rails.logger.warn e
        end
        Process.wait(pid)
      end

      yield $CHILD_STATUS.inspect if block_given? && !success?($CHILD_STATUS)

      force_fail ? false : $CHILD_STATUS
    end

    sig { params(child_status: T.untyped).returns(T.nilable(T::Boolean)) }
    def success?(child_status)
      child_status&.exitstatus&.zero?
    end
  end
end
